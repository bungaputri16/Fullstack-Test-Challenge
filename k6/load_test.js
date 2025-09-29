import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 100,          // 100 Virtual Users
    duration: '10s',   // Jalankan 10 detik
    rps: 1000          // Target 1000 request/sec
};

// Hostname gateway-service sesuai docker network
const BASE_URL = 'http://gateway-service:5000';

export default function () {
    const payload = JSON.stringify({
        productId: 11,
        qty: 1
    });

    const params = {
        headers: { 'Content-Type': 'application/json' }
    };

    const res = http.post(`${BASE_URL}/orders`, payload, params);

    check(res, { 'status is 200': (r) => r.status === 200 });


}
