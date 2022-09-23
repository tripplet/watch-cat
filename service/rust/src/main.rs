use std::{error, sync::Arc, thread, time::Duration, net::{ToSocketAddrs, SocketAddr}};

// Logging
use log::{error, info, LevelFilter};
use simple_logger::SimpleLogger;

// Network related
use ureq::{Agent, Response};
use url::Url;

// Parse cmdline arguments
use clap::Parser;

// The main config
#[derive(Debug, Parser)]
#[clap(
    version,
    about = "Service for sending requests to the watchcat backend."
)]
struct Config {
    /// Url where to send requests
    #[clap(short, parse(try_from_str = Url::parse), long, env)]
    url: Url,

    /// HTTP method to use
    #[clap(long, default_value = "POST", env)]
    method: String,

    /// Secret key to use
    #[clap(long, env, hide_env_values = true)]
    key: Option<String>,

    /// Repeat request after interval, with units 'ms', 's', 'm', 'h', e.g. 2m30s
    #[clap(long, parse(try_from_str = humantime::parse_duration), env)]
    repeat: Option<Duration>,

    /// Timeout for http request in seconds
    #[clap(long, default_value = "30", env)]
    timeout: u16,

    /// Check dns every x seconds before first request, for faster inital signal in case of long allowed timeout
    #[clap(long, env)]
    checkdns: Option<u16>,

    /// Do not send uptime in heartbeat requests
    #[clap(long, env)]
    nouptime: bool,

    /// Verbose mode
    #[clap(short, long)]
    verbose: bool,
}

fn main() {
    let cfg = Config::parse(); // Parse arguments

    // Initialize logger
    SimpleLogger::new().init().unwrap();
    if cfg.verbose {
        log::set_max_level(LevelFilter::Info);
    } else {
        log::set_max_level(LevelFilter::Warn);
    }

    info!("{:?}", &cfg);

    // Create HTTP agent
    let agent = ureq::builder()
        .user_agent(&*format!("watchcat-service/{}", env!("CARGO_PKG_VERSION")))
        .timeout(Duration::from_secs(cfg.timeout.into()))
        .tls_connector(Arc::new(native_tls::TlsConnector::new().unwrap()))
        .build();

    // Preform DNS pre check loop if configured
    match cfg.checkdns {
        Some(interval) if interval > 0 => check_dns(interval, cfg.url.domain().unwrap()),
        _ => (),
    }

    loop {
        // Send the HTTP request
        let resp = send_request(&agent, &cfg);
        match resp {
            Ok(resp) => {
                let status = resp.status();
                match resp.into_string() {
                    Ok(resp_str) => info!("Response {status}: \"{resp_str}\""),
                    Err(err) => error!("Error receiving response: {err:?}"),
                }
            }
            Err(err) => error!("Error during request: {err}"),
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

    let socket_addr = format!("{}:443", domain);

    loop {
        match socket_addr.to_socket_addrs() {
            Ok(mut dns_entries) => {
                info!(
                    "DNS lookup successful: {domain} -> {ips}",
                    ips = to_ip_list(&mut dns_entries)
                );
                break;
            }
            Err(err) => {
                info!("DNS lookup failed '{err}', retrying in {interval} seconds");
                thread::sleep(Duration::from_secs(interval.into()));
            }
        }
    }
}

/// Send a HTTP request
fn send_request(agent: &Agent, cfg: &Config) -> Result<Response, Box<dyn error::Error>> {
    let mut request = agent.request(cfg.method.as_str(), cfg.url.as_str());

    // Add current uptime as query parameter
    if !cfg.nouptime {
        request = request.query("uptime", uptime_lib::get()?.as_secs().to_string().as_str());
    }

    // Add the secret key as authorization header
    if let Some(key) = &cfg.key {
        request = request.set("Authorization", format!("Bearer {key}").as_str());
    }

    request = request.set("Content-Length", "0");

    Ok(request.call()?)
}

fn to_ip_list(ips: &mut dyn Iterator<Item = SocketAddr>) -> String {
    ips
        .map(|addr| addr.ip().to_string())
        .reduce(|a, b| format!("{a}, {b}")).unwrap_or_else(|| "".to_string())
}
