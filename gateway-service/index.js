const express = require('express');
const axios = require('axios');
const { v4: uuidv4 } = require('uuid');

const app = express();
app.use(express.json());

// URL service lain (Docker internal hostnames)
const PRODUCT_SERVICE_URL = 'http://product-service:3000';
const ORDER_SERVICE_URL = 'http://order-service:4000';

// Middleware: generate request ID
app.use((req, res, next) => {
  const requestId = uuidv4();
  req.requestId = requestId;
  res.setHeader('X-Request-ID', requestId);
  console.log(`[${requestId}] ${req.method} ${req.url}`);
  next();
});

// Combined endpoint: order + product
app.get('/orders/details/:id', async (req, res) => {
    const { id } = req.params;

    try {
        // Ambil order berdasarkan productId
        const orderRes = await axios.get(`${ORDER_SERVICE_URL}/orders/product/${id}`);
        const orders = orderRes.data;

        // Ambil detail product
        const productRes = await axios.get(`${PRODUCT_SERVICE_URL}/products/${id}`);
        const product = productRes.data;

        if (!product) {
          return res.status(404).json({ error: 'Product not found' });
        }

        res.json({ product, orders });
    } catch (err) {
        console.error(`[${req.requestId}] Gateway error:`, err.message);
        if (err.response && err.response.status === 404) {
            return res.status(404).json({ error: 'Data not found', details: err.response.data });
        }
        res.status(500).json({ error: 'Failed to fetch order details', details: err.message });
    }
});

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok', requestId: req.requestId });
});

// Jalankan server
const PORT = 5000;
app.listen(PORT, () => console.log(`Gateway service running on port ${PORT}`));
