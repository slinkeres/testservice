async function getOrder() {
    const orderId = document.getElementById('orderId').value.trim();
    if (!orderId) {
        showError('Please enter an Order ID');
        return;
    }

    showLoading();
    
    try {
        const response = await fetch(`http://localhost:8080/order/${orderId}`);
        
        if (response.ok) {
            const order = await response.json();
            displayOrder(order);
        } else if (response.status === 404) {
            showError('Order not found');
        } else {
            showError('Error fetching order data');
        }
    } catch (error) {
        showError('Error: ' + error.message);
    }
}

function displayOrder(order) {
    const resultDiv = document.getElementById('result');
    

    const dateCreated = new Date(order.date_created).toLocaleString();
    let html = `
        <div class="order-detail">
            <h3>Order Information</h3>
            <div class="detail-row">
                <div class="detail-label">Order UID:</div>
                <div class="detail-value">${order.order_uid}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Track Number:</div>
                <div class="detail-value">${order.track_number}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Entry:</div>
                <div class="detail-value">${order.entry}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Locale:</div>
                <div class="detail-value">${order.locale}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Customer ID:</div>
                <div class="detail-value">${order.customer_id}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Delivery Service:</div>
                <div class="detail-value">${order.delivery_service}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Date Created:</div>
                <div class="detail-value">${dateCreated}</div>
            </div>
        </div>
        
        <div class="order-detail">
            <h3>Delivery Information</h3>
            <div class="detail-row">
                <div class="detail-label">Name:</div>
                <div class="detail-value">${order.delivery.name}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Phone:</div>
                <div class="detail-value">${order.delivery.phone}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Zip:</div>
                <div class="detail-value">${order.delivery.zip}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">City:</div>
                <div class="detail-value">${order.delivery.city}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Address:</div>
                <div class="detail-value">${order.delivery.address}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Region:</div>
                <div class="detail-value">${order.delivery.region}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Email:</div>
                <div class="detail-value">${order.delivery.email}</div>
            </div>
        </div>
        
        <div class="order-detail">
            <h3>Payment Information</h3>
            <div class="detail-row">
                <div class="detail-label">Transaction:</div>
                <div class="detail-value">${order.payment.transaction}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Currency:</div>
                <div class="detail-value">${order.payment.currency}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Provider:</div>
                <div class="detail-value">${order.payment.provider}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Amount:</div>
                <div class="detail-value">${order.payment.amount}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Bank:</div>
                <div class="detail-value">${order.payment.bank}</div>
            </div>
        </div>
        
        <div class="order-detail">
            <h3>Items (${order.items.length})</h3>
            <table class="items-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Brand</th>
                        <th>Price</th>
                        <th>Sale</th>
                        <th>Total Price</th>
                        <th>Status</th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    order.items.forEach(item => {
        html += `
            <tr>
                <td>${item.name}</td>
                <td>${item.brand}</td>
                <td>${item.price}</td>
                <td>${item.sale}%</td>
                <td>${item.total_price}</td>
                <td>${item.status}</td>
            </tr>
        `;
    });
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    resultDiv.innerHTML = html;
}

function showError(message) {
    const resultDiv = document.getElementById('result');
    resultDiv.innerHTML = `<div class="error">${message}</div>`;
}

function showLoading() {
    const resultDiv = document.getElementById('result');
    resultDiv.innerHTML = '<div class="loading">Loading...</div>';
}


document.getElementById('orderId').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        getOrder();
    }
});