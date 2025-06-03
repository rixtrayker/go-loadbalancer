import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

const errorRate = new Rate("errors");

export const options = {
  scenarios: {
    spike: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "10s", target: 50 }, // Normal load
        { duration: "1m", target: 50 }, // Stay at normal load
        { duration: "10s", target: 500 }, // Spike to 500 users
        { duration: "3m", target: 500 }, // Stay at spike
        { duration: "10s", target: 50 }, // Back to normal
        { duration: "3m", target: 50 }, // Stay at normal
        { duration: "10s", target: 0 }, // Ramp down
      ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<2000"], // 95% of requests should be below 2s
    http_req_failed: ["rate<0.1"], // Less than 10% of requests should fail
    errors: ["rate<0.15"], // Less than 15% of requests should error
  },
};

export default function () {
  const response = http.get("http://localhost:8080/", {
    tags: { name: "spike-test" },
  });

  const success = check(response, {
    "is status 200": (r) => r.status === 200,
    "has correct content": (r) => r.body.includes("Backend"),
  });

  errorRate.add(!success);

  sleep(1);
}
