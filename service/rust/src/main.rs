use std::{
    borrow::Cow,
    error,
    net::{SocketAddr, ToSocketAddrs},
    str::FromStr,
    thread,
    time::Duration,
};

use chrono::{TimeDelta, Utc};
use log::{error, info, LevelFilter};
use reqwest::{
    blocking::{Client, Response},
    Method, Url,
};
use simple_logger::SimpleLogger;

type Result<T> = std::result::Result<T, Box<dyn error::Error>>;

#[derive(Debug, Copy, Clone, PartialEq, Eq)]
enum UptimeMode {
    UpTime,
    BootTime,
    None,
}

// The main config
#[derive(Debug, clap::Parser)]
#[clap(version, about = "Service for sending requests to the watchcat backend.")]
struct Config {
    /// Url where to send requests
    #[clap(short, long, env)]
    url: Url,

    /// HTTP method to use
    #[clap(long, default_value = "POST", env)]
    method: Method,

    /// Secret key to use
    #[clap(long, env, hide_env_values = true)]
    key: Option<String>,

    /// Repeat request after interval, with units 'ms', 's', 'm', 'h', e.g. 2m30s
    #[clap(long, value_parser = humantime::parse_duration, env)]
    repeat: Option<Duration>,

    /// Timeout for http request in seconds
    #[clap(long, default_value = "30", env)]
    timeout: u16,

    /// Check dns every x seconds before first request, for faster initial signal in case of long allowed timeout
    #[clap(long, env)]
    checkdns: Option<u16>,

    /// Send uptime in heartbeat requests
    #[clap(long, env)]
    uptime: UptimeMode,

    /// Verbose mode
    #[clap(short, long)]
    verbose: bool,
}

fn main() {
    let cfg: Config = clap::Parser::parse(); // Parse arguments

    // Initialize logger
    SimpleLogger::new().init().unwrap();
    if cfg.verbose {
        log::set_max_level(LevelFilter::Info);
    } else {
        log::set_max_level(LevelFilter::Warn);
    }

    info!("{:?}", &cfg);

    // Preform DNS pre check loop if configured
    match cfg.checkdns {
        Some(interval) if interval > 0 => check_dns(interval, cfg.url.domain().unwrap()),
        _ => (),
    }

    loop {
        // Create HTTP agent
        let client = Client::builder()
            .user_agent(format!("watchcat-service/{}", env!("CARGO_PKG_VERSION")))
            .timeout(Duration::from_secs(cfg.timeout.into()))
            .build()
            .unwrap();

        // Send the HTTP request
        let resp = send_request(&client, &cfg);
        match resp {
            Ok(resp) => {
                let status = resp.status();
                match resp.text() {
                    Ok(resp_str) => info!("Response {status}: \"{resp_str}\""),
                    Err(err) => error!("Error receiving response: {err:?}"),
                }
            }
            Err(err) => error!("Error during request: {err}"),
        }

        if let Some(repeat_interval) = cfg.repeat {
            // Save memory until next request
            drop(client);
            thread::sleep(repeat_interval);
        } else {
            break;
        }
    }
}

/// Try to resolve the domains DNS entry in a loop until it is successful.
fn check_dns(interval: u16, domain: &str) {
    info!("Checking DNS");

    let socket_addr = format!("{domain}:443");

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
fn send_request(agent: &Client, cfg: &Config) -> Result<Response> {
    let mut request = agent.request(cfg.method.clone(), cfg.url.as_str());

    // Add current uptime
    if let Some(info) = cfg.uptime.to_query()? {
        request = request.query(&[info]);
    }

    // Add the secret key as authorization header
    if let Some(key) = &cfg.key {
        request = request.header("Authorization", format!("Bearer {key}").as_str());
    }

    request = request.header("Content-Length", "0");

    info!("Request: {request:?}");
    Ok(request.send()?)
}

fn to_ip_list(ips: &mut dyn Iterator<Item = SocketAddr>) -> String {
    ips.map(|addr| addr.ip().to_string())
        .reduce(|addr1, addr2| format!("{addr1}, {addr2}"))
        .unwrap_or_default()
}

impl FromStr for UptimeMode {
    type Err = String;

    fn from_str(s: &str) -> std::result::Result<Self, Self::Err> {
        match s.to_uppercase().as_str() {
            "NONE" => Ok(UptimeMode::None),
            "UPTIME" => Ok(UptimeMode::UpTime),
            "BOOTTIME" | "BOOT" => Ok(UptimeMode::BootTime),
            _ => Err(format!("Invalid value: {s} for UptimeMode")),
        }
    }
}

impl UptimeMode {
    fn to_query(self) -> Result<Option<(Cow<'static, str>, String)>> {
        match self {
            UptimeMode::UpTime => {
                let uptime = uptime_lib::get()?;
                Ok(Some((Cow::Borrowed("uptime"), uptime.as_secs().to_string())))
            }
            UptimeMode::BootTime => {
                let boot = Utc::now()
                    .checked_sub_signed(TimeDelta::from_std(uptime_lib::get()?)?)
                    .ok_or("Unable to get boot time")?;

                Ok(Some((Cow::Borrowed("boot"), boot.to_rfc3339())))
            }
            UptimeMode::None => Ok(None),
        }
    }
}
