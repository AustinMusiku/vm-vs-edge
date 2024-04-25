# HTTP Service Response Time Experiment
## Introduction
This repository contains an experiment aimed at comparing the response time (RT) performance between a Golang HTTP service hosted on a Virtual Private Server (VPS) and one deployed on a serverless edge platform (Fastly). The experiment was conducted from Nairobi, Kenya, with the VPS located in Paris, France, and the nearest Fastly edge node located in Johannesburg, South Africa.
I had initially planned to use cloudflare's workers platform for the edge server since it has edge nodes in Nairobi, 
but it was a bit challenging to make it work with go. Hopefully, I will succeed in the future and prepare a sequel to this experiment.

## Setup
- VPS Server Code and Configuration: The code and configuration for the VPS-based HTTP service are located in the `./vm` directory.
- Edge Server Code and Configuration: The code and configuration for the serverless edge platform-based HTTP service are located in the `./edge` directory.
- Experiment Entry Point: The entry point for the experiment is a CLI application, located at `./benchmark.go`. This CLI application accepts one flag, `n`, and two arguments, `vps-url` and `edge-url`. The application sends `n` HTTP requests to each of the services and records the `mean`, `minimum`, and `maximum` response time.

## Usage

To conduct the experiment with custom parameters, follow these steps:

1. Clone the repository: `git clone https://github.com/AustinMusiku/vm-vs-edge`
2. Navigate to the project directory: `cd vm-vs-edge`
3. Deploy the VPS and edge server applications as per the instructions in the respective directories.
4. Compile the CLI application:
    ```sh
    go build -o=./bin/benchmark benchmark.go
    ```
5. Run the CLI application with custom parameters:
    ```sh
    ./bin/benchmark -n <number-of-requests> <edge-url> <vps-url>
    ```

## Results
After running the CLI application with the default setting of `n=100`, the following results were obtained:

| Service Type | Mean RT (ms) | Min RT (ms) | Max RT (ms) |
|--------------|--------------|-------------|-------------|
| Fastly Edge Platform | 246.263 | 213.185 | 295.132 |
| Virtual Private Server (VPS) | 451.929 | 428.668 | 495.378 |


## Summary
The experiment demonstrates that the response time performance of the HTTP service hosted on the Fastly edge platform outperforms the one hosted on the Virtual Private Server. The average RT for the Fastly edge platform was approximately `45.5%` faster than that of the VPS. These findings suggest that leveraging edge platforms like Fastly compute or cloudflare workers can significantly improve response time for HTTP services, especially for users located in regions closer to edge nodes.

## Caveats
This experiment was conducted with a limited number of requests and in a specific geographical context (Nairobi, Kenya). The results may vary under different conditions, such as varying the number of requests, geographical locations, network conditions, and the efficiency of the underlying infrastructure. The specific configuration of the VPS and edge platform, such as server specifications, caching mechanisms, and network routing, could also impact the response time performance.

## Conclusion
There exists a potential significant performance benefit in using edge platforms like Fastly for hosting HTTP services compared to traditional VPS hosting. The distributed nature of edge computing allows for lower latency and faster response times, particularly for users located closer to edge nodes. This performance advantage can be crucial for applications that require low-latency responses and aim to enhance user experience.

It's also worth mentioning that there exists a trade-off between the performance benefits of edge computing and the cost implications of using edge platforms compared to traditional VPS hosting. You need to consider factors such as scalability, reliability, and cost-effectiveness when choosing between VPS and edge platforms for hosting their web services.
