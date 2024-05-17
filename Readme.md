# HTTP Service Response Time Experiment
## Introduction
This repository contains an experiment aimed at comparing the response time (RT) performance between a Golang HTTP service hosted on a Virtual Private Server (VPS) and one deployed on a serverless edge platform (Fastly). The experiment was conducted from Nairobi, Kenya, with the VPS located in Paris, France, and the nearest Fastly edge node located in Johannesburg, South Africa.
I had initially planned to use cloudflare's workers platform for the edge server since it has edge nodes in Nairobi, 
but it was a bit challenging to make it work with go. Hopefully, I will succeed in the future and prepare a sequel to this experiment.

## Setup
### HTTP Services
The two services expose two endpoints each: `/` and `/characters`. The `/` endpoint handler iterates over a list of names, `characters`, and appends each name to a html ordered list template string. The final html list element is sent over in a html page. The `/characters` endpoint just spits out the list of characters in JSON format.

### CLI Application
The CLI app automates the process of sending HTTP requests to the two services and recording the response time metrics. The response time metrics are then used to compute the relevant numerical summaries.

The app accepts the following arguments:
1. `-n`: The number of requests to send to each service.
2. `edge-url`: The URL of the edge server.
3. `vps-url`: The URL of the VPS server.


### Folder structure
- VPS Server Code and Configuration: The code and configuration for the VPS-based HTTP service are located in the `./vm` directory.
- Edge Server Code and Configuration: The code and configuration for the serverless edge platform-based HTTP service are located in the `./edge` directory.
- Experiment Entry Point: The entry point for the experiment is a CLI application, located at `./benchmark.go`. This CLI application accepts one flag, `n`, and two arguments, `vps-url` and `edge-url`. The application sends `n` HTTP requests to each of the services and records the `mean`, `p25`, `p50`, `p75`, `p90`, `minimum`, and `maximum` response time for each service.

## Usage

To conduct the experiment with custom parameters, follow these steps:

1. Clone the repository: `git clone https://github.com/AustinMusiku/vm-vs-edge`
2. Navigate to the project directory: `cd vm-vs-edge`
3. Deploy the VPS and edge server applications as per the instructions in the respective directories.
4. Compile the CLI application:
    ```sh
    go build -o=./bin/benchmark benchmark.go
    ```
    or using make:
    ```sh
    make compile
    ```
5. Run the CLI application with custom parameters:
    ```sh
    ./bin/benchmark -n <number-of-requests> <edge-url> <vps-url>
    ```
    or using make:
    ```sh
    make run
    ```

## Results
After running the CLI application with the default setting of `n=100`, the following results were obtained:

| RT Metric(ms) | Fastly Edge Platform | Virtual Private Server (VPS) |
|--------------|--------------|--------------|
| Min | 127.786 | 315.221 |
| p25 | 153.289 | 350.572 |
| p50 | 161.830 | 365.363 |
| p75 | 183.662 | 390.331 |
| p90 | 194.277 | 406.690 |
| Avg | 171.251 | 370.293 |
| Max | 296.074 | 424.539 | 


## Summary
The experiment demonstrates that the response time performance of the HTTP service hosted on the Fastly edge platform outperforms the one hosted on the Virtual Private Server. The median response time for the Fastly edge platform was approximately `55.71%` faster than that of the VPS. These findings suggest that leveraging edge platforms like Fastly compute or cloudflare workers can significantly improve response time for HTTP services, especially for users located in regions closer to edge nodes.

## Caveats
This experiment was conducted with a limited number of requests and in a specific geographical context (Nairobi, Kenya). The results may vary under different conditions, such as varying the number of requests, geographical locations, network conditions, and the efficiency of the underlying infrastructure. The specific configuration of the VPS and edge platform, such as server specifications, caching mechanisms, and network routing, could also impact the response time performance.

## Conclusion
There exists a potential significant performance benefit in using edge platforms like Fastly for hosting HTTP services compared to traditional VPS hosting. The distributed nature of edge computing allows for lower latency and faster response times, particularly for users located closer to edge nodes. This performance advantage can be crucial for applications that require low-latency responses and aim to enhance user experience.

It's also worth mentioning that there exists a trade-off between the performance benefits of edge computing and the cost implications of using edge platforms compared to traditional VPS hosting. You need to consider factors such as scalability, reliability, and cost-effectiveness when choosing where to host your web services.
