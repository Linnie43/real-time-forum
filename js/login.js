// login.js
document.getElementById('login-form').addEventListener('submit', async (event) => {
    event.preventDefault();

    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const loginData = {
        username: username,
        password: password
    };

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });

        if (response.ok) {
            const data = await response.json();
            alert(data.message); // Show success message
            // Redirect to feed or home page
            window.location.href = '/feed'; 
        } else {
            const errorData = await response.json();
            alert(errorData.message); // Show error message
        }
    } catch (error) {
        alert('An error occurred while logging in');
    } finally {
        // restore the login button
        loginButton.disabled = false;
        loginButton.textContent = 'Login';
    }
});
