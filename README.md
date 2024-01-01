# SlowSloth - Slow HTTP Attack Simulator

![slowsloth-logo](https://github.com/mduygu/slowsloth/assets/61627415/90cf459e-77a6-439e-9759-317459e652e1)

SlowSloth is a Go application designed to simulate slow HTTP attacks on APIs, helping you perform stress tests on your web services. This README provides an overview of the application and explains how to use it effectively.

## Release
Current Release: v1.1.1

## Table of Contents

- [Overview](#overview)
 - [Key Features](#key-features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Overview

SlowSloth is a powerful Go application designed to facilitate stress testing of web services and APIs by simulating slow HTTP attacks. It provides a flexible and versatile environment for assessing the robustness and resilience of your services when faced with slow client connections.

### Key Features

- **Slow HTTP Attack Simulation:** SlowSloth allows you to simulate slow HTTP requests, emulating the behavior of clients with sluggish network conditions. This enables you to evaluate how well your service can handle slow or delayed requests.

- **Customizable Parameters:** You have full control over the testing parameters. You can specify the target URL, choose between HTTP methods (GET or POST), provide request data for POST requests, set the concurrency level to simulate multiple clients, and control the delay between requests.

- **Real-time Monitoring:** SlowSloth provides real-time monitoring of the testing process. It tracks the service's availability and active connections, giving you insights into how your service is performing during the test.

- **Concurrent Testing:** With the ability to specify the concurrency level, you can simulate multiple clients simultaneously, putting a substantial load on your service to assess its scalability.

- **User-Agent Rotation:** SlowSloth rotates through a set of User-Agent strings to mimic different client devices, providing a more realistic testing scenario.

- **Easy-to-Use:** SlowSloth's simple command-line interface makes it easy to set up and run tests quickly, allowing you to focus on evaluating your service's performance.

Whether you want to stress test a production API, evaluate the performance of a new service, or identify potential bottlenecks in your web application, SlowSloth is a valuable tool for ensuring the reliability and stability of your online services.


## Getting Started

### Prerequisites

Before using SlowSloth, ensure you have the following prerequisites installed on your system:

- [Go](https://golang.org/dl/): You must have Go installed to build and run the application.
- [Git](https://git-scm.com/downloads): Git is required to clone the repository (if not already done).

### Installation

1. Clone the SlowSloth repository to your local machine:

   ```sh
   git clone https://github.com/mduygu/slowsloth.git
   
## Usage

Run the application with the desired parameters:
   ```sh
   ./SlowSloth -u <Target_URL> -m <HTTP_Method> -d <Request_Data> -c <Concurrency_Level> -delay <Delay_in_seconds>
   ```
Replace the placeholders with your specific values:
- **<Target_URL>**: The URL of the target service you want to test.
- **<HTTP_Method>**: The HTTP method to use (GET or POST).
- **<Request_Data>**: Data to include in the POST request (if applicable).
- **<Concurrency_Level>**: Number of concurrent requests to send.
- **<Delay_in_seconds>**: Delay in seconds between sending header data.
- For example, to test a target URL "http://example.com" with 10 concurrent requests using the - GET method and a 5-second delay, you can run:
  
   ```sh
   ./SlowSloth -u http://example.com -m GET -c 10 -delay 5
   ```

## Contributing

Contributions are welcome! If you'd like to contribute to SlowSloth, please follow these guidelines:

- Fork the repository on GitHub.
- Clone your forked repository to your local machine.
- Create a new branch for your feature or bug fix.
- Make your changes and commit them with descriptive commit messages.
- Push your changes to your forked repository.
- Submit a pull request to the main SlowSloth repository.
