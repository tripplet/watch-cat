use std::{error::Error, thread, time::Duration};

// Logging
use log::{error, info, LevelFilter};
use simple_logger::SimpleLogger;

// Network related
use dns_lookup;
use http::{header, HeaderMap, HeaderValue, Method};
use reqwest::blocking::{Client, Response};
use url::Url;

// Other stuff
use humantime;
use structopt::StructOpt;
use uptime_lib;

// The main config
#[derive(Debug, StructOpt)]
#[structopt(about = "Service for sending requests to the watchcat backend.")]
struct Config {
    /// Url where to send requests
    #[structopt(short, parse(try_from_str = Url::parse), long, env)]
    url: Url,

    /// HTTP method to use
    #[structopt(long, default_value = "POST", env)]
    method: String,

    /// Secret key to use
    #[structopt(long, env, hide_env_values = true)]
    key: Option<String>,

    /// Repeat request after interval, with units 'ms', 's', 'm', 'h', e.g. 2m30s
    #[structopt(long, parse(try_from_str = humantime::parse_duration), env)]
    repeat: Option<Duration>,

    /// Timeout for http request in seconds
    #[structopt(long, env)]
    timeout: Option<u16>,

    /// Check dns every x seconds before first request, for faster inital signal in case of long allowed timeout
    #[structopt(long, env)]
    checkdns: Option<u16>,

    /// Do not send uptime in heartbeat requests
    #[structopt(long, env)]
    nouptime: bool,

    /// Verbose mode
    #[structopt(short, long)]
    verbose: bool,
}

fn main() {
    let cfg = Config::from_args(); // Parse arguments

    // Initialize logger
    SimpleLogger::new().init().unwrap();
    if cfg.verbose {
        log::set_max_level(LevelFilter::Info);
    } else {
        log::set_max_level(LevelFilter::Warn);
    }

    info!("{:?}", &cfg);

    // Create HTTP client
    let mut client_builder = Client::builder();
    let mut headers = HeaderMap::new();

    // Add timeout if specified
    if let Some(http_timeout) = cfg.timeout {
        client_builder = client_builder.timeout(Duration::from_secs(http_timeout.into()));
    }

    // Add the secret key as authorization header
    if let Some(key) = &cfg.key {
        headers.insert(header::AUTHORIZATION, HeaderValue::from_str(&key).unwrap());
    }

    headers.insert(header::CONTENT_LENGTH, HeaderValue::from_static("0"));
    client_builder = client_builder.default_headers(headers);
    let client = client_builder.build().unwrap();

    // Preform DNS pre check loop if configured
    if let Some(check_dns_interval) = cfg.checkdns {
        if check_dns_interval > 0 {
            check_dns(check_dns_interval, cfg.url.domain().unwrap());
        }
    }

    loop {
        // Send the HTTP request
        match send_request(&client, &cfg) {
            Ok(resp) => {
                info!(
                    "Response {}: \"{}\"",
                    resp.status().as_u16(),
                    resp.text().unwrap_or_default()
                );
            }
            Err(err) => {
                error!("Error during request: {}", err);
            }
        }

        if let Some(repeat_interval) = cfg.repeat {
            thread::sleep(repeat_interval);
        } else {
            break;
        }
    }
}

/// Try to resolve the domains DNS entry in a loop until it is succesful.
fn check_dns(interval: u16, domain: &str) {
    info!("Checking DNS");

    loop {
        if let Ok(dns_entries) = dns_lookup::lookup_host(domain) {
            info!("DNS lookup successful: {} -> {}", domain, dns_entries[0]);
            break;
        } else {
            info!("DNS lookup failed, retrying in {} seconds", interval);
            thread::sleep(Duration::from_secs(interval.into()));
        }
    }
}

/// Send a HTTP request
fn send_request(client: &Client, cfg: &Config) -> Result<Response, Box<dyn Error>> {
    let mut request_builder =
        client.request(Method::from_bytes(cfg.method.as_bytes())?, cfg.url.as_str());

    // Add current uptime as query parameter
    if !cfg.nouptime {
        request_builder =
            request_builder.query(&[("uptime", uptime_lib::get()?.as_secs().to_string())]);
    }

    let request = request_builder.build()?;

    info!("Requesting {}: {}", cfg.method, request.url().as_str());
    Ok(client.execute(request)?)
}
