import http from "k6/http";
import { check, sleep } from "k6";

const TARGET_URL = __ENV.TARGET_URL || "http://localhost:8080";
const TEST_DURATION = __ENV.TEST_DURATION || "1m";
const VUS = parseInt(__ENV.VUS || "20");

export const options = {
  stages: [
    { duration: "30s", target: VUS }, // Ramp up to target VUs
    { duration: TEST_DURATION, target: VUS }, // Stay at target VUs
    { duration: "30s", target: 0 }, // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"], // 95% of requests should be below 500ms
    http_req_failed: ["rate<0.01"], // Less than 1% of requests should fail
  },
};

export default function () {
  const response = http.get(TARGET_URL);

  check(response, {
    "is status 200": (r) => r.status === 200,
    "has correct content": (r) => r.body.includes("Backend"),
  });

  sleep(1);
}
