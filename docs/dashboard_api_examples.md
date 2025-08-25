# Dashboard API Endpoints

## GET /api/admin/dashboard/stats

Retrieve dashboard statistics for a specified time period.

### Authentication

This endpoint requires admin authentication.

### Query Parameters

| Parameter    | Type   | Required | Description                                                  |
| ------------ | ------ | -------- | ------------------------------------------------------------ |
| `start_date` | string | No       | Start date in YYYY-MM-DD format                              |
| `end_date`   | string | No       | End date in YYYY-MM-DD format                                |
| `days`       | int    | No       | Number of days from current date (alternative to date range) |

If no parameters are provided, defaults to the last 30 days.

### Example Requests

#### Get stats for the last 30 days (default)

```http
GET /api/admin/dashboard/stats
```

#### Get stats for a specific number of days

```http
GET /api/admin/dashboard/stats?days=7
```

#### Get stats for a specific date range

```http
GET /api/admin/dashboard/stats?start_date=2025-01-01&end_date=2025-01-31
```

### Response

```json
{
  "success": true,
  "message": "Dashboard statistics retrieved successfully",
  "data": {
    "total_revenue": 4567890,
    "total_orders": 234,
    "total_customers": 1247,
    "new_customers": 23,
    "total_products": 156,
    "low_stock_products": 8,
    "revenue_change": {
      "value": 15.5,
      "direction": "up"
    },
    "orders_change": {
      "value": 8.2,
      "direction": "up"
    },
    "recent_orders": [
      {
        "id": 1001,
        "order_number": "ORD-2025-001",
        "customer_name": "John Doe",
        "customer_email": "john@example.com",
        "total_amount": 12345,
        "status": "completed",
        "created_at": "2025-08-20T10:30:00Z"
      }
    ],
    "top_products": [
      {
        "product_id": 1,
        "product_name": "Wireless Headphones",
        "variant_id": 1,
        "variant_name": "Black",
        "quantity_sold": 45,
        "revenue": 225000
      }
    ],
    "period_start": "2025-07-22T00:00:00Z",
    "period_end": "2025-08-21T23:59:59Z"
  }
}
```

### Response Fields

- `total_revenue`: Total revenue in cents for the specified period
- `total_orders`: Total number of orders placed in the period
- `total_customers`: Total number of registered customers (all time)
- `new_customers`: Number of new customers registered in the period
- `total_products`: Total number of active products in the system
- `low_stock_products`: Number of products with stock at or below the low stock threshold (10 units)
- `revenue_change`: Percentage change in revenue compared to the previous equivalent period
  - `value`: Absolute percentage change (e.g., 15.5 for 15.5% change)
  - `direction`: "up", "down", or "stable"
- `orders_change`: Percentage change in orders compared to the previous equivalent period
  - `value`: Absolute percentage change
  - `direction`: "up", "down", or "stable"
- `recent_orders`: Array of recent orders (limited to 10) with basic information
- `top_products`: Array of top-selling products (limited to 10) with sales data
- `period_start`: Start of the queried period
- `period_end`: End of the queried period

### Error Responses

#### 400 Bad Request

```json
{
  "success": false,
  "error": "Invalid start_date format. Use YYYY-MM-DD"
}
```

#### 401 Unauthorized

```json
{
  "success": false,
  "error": "Authentication required"
}
```

#### 403 Forbidden

```json
{
  "success": false,
  "error": "Admin access required"
}
```

#### 500 Internal Server Error

```json
{
  "success": false,
  "error": "Failed to retrieve dashboard statistics"
}
```

### Notes

- All monetary values are returned in cents (e.g., $123.45 = 12345)
- Revenue calculations include only paid, shipped, or completed orders
- Percentage changes compare the current period with the previous equivalent period
  - For a 30-day period, it compares with the previous 30 days
  - If there's no previous period data, changes show as 100% "up" if current period has data, or 0% "stable" if both periods have no data
- Low stock threshold is set to 10 units or less
- Top products are ranked by quantity sold
- Recent orders are ordered by creation date (most recent first)
- Guest orders are included with customer_name as "Guest"
