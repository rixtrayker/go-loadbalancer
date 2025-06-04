import http from "k6/http";
import { check, sleep } from "k6";
import { Rate } from "k6/metrics";

const errorRate = new Rate("errors");

export const options = {
  scenarios: {
    stress: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "2m", target: 100 }, // Ramp up to 100 users
        { duration: "5m", target: 100 }, // Stay at 100 users
        { duration: "2m", target: 200 }, // Ramp up to 200 users
        { duration: "5m", target: 200 }, // Stay at 200 users
        { duration: "2m", target: 0 }, // Ramp down to 0 users
      ],
      gracefulRampDown: "30s",
    },
  },
  thresholds: {
    http_req_duration: ["p(95)<1000"], // 95% of requests should be below 1s
    http_req_failed: ["rate<0.05"], // Less than 5% of requests should fail
    errors: ["rate<0.1"], // Less than 10% of requests should error
  },
};

export default function () {
  const response = http.get("http://localhost:8080/", {
    tags: { name: "stress-test" },
  });

  const success = check(response, {
    "is status 200": (r) => r.status === 200,
    "has correct content": (r) => r.body.includes("Backend"),
  });

  errorRate.add(!success);

  sleep(1);
}
