document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior

            // Gather form data
            const username = document.querySelector('input[name="username"]').value;
            const password = document.querySelector('input[name="password"]').value;

            // Regular expressions for validation
            const usernameRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const passwordRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;

            // Perform validation
            if (!usernameRegex.test(username)) {
                alert('Invalid username. Use only Latin characters and special symbols.');
                return;
            }

            if (!passwordRegex.test(password)) {
                alert('Invalid password. Use only Latin characters and special symbols.');
                return;
            }
            if (password.length<8) {
                alert('Password should contain at least 8 characters');
                return;
            }

            // Create an object with the form data
            const formData = {
                username,
                password,
            };

            // Send a POST request to the server
            try {
                const response = await fetch('/auth/sign-in', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });


                if (response.ok) {
                    // Registration was successful
                    const data = await response.json();
                    const jwtToken = data.token;
                    localStorage.setItem('token',jwtToken);
                    console.log('Successful');
                    if (username==='admin'){
                        window.location.assign('/admin/user/page')
                    }else{
                        window.location.assign('/login')
                    }

                } else {
                    // Registration failed
                    console.error('Failed');
                    // Handle the failure case (e.g., show an error message to the user).
                }
                    // You can handle the success case here (e.g., show a success message or redirect the user).

            } catch (error) {
                console.error('Error:', error);
                // Handle network or other errors here.
            }
        });
    }
});
