document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior

            // Gather form data
            const name = document.querySelector('input[name="name"]').value;
            const username = document.querySelector('input[name="username"]').value;
            const email = document.querySelector('input[name="email"]').value;
            const password = document.querySelector('input[name="password"]').value;
            const confirmPassword = document.querySelector('input[name="confirmPassword"]').value;


            const usernameRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const passwordRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const confirmPasswordRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const emailRegex=/^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$/
            // Perform validation
            if (!usernameRegex.test(username)) {
                alert('Invalid username. Use only Latin characters and special symbols.');
                return;
            }
            if (!passwordRegex.test(password)) {
                alert('Invalid password. Use only Latin characters and special symbols.');
                return;
            }
            if (!confirmPasswordRegex.test(password)) {
                alert('Invalid confirm password. Use only Latin characters and special symbols.');
                return;
            }
            if (!emailRegex.test(email)) {
                alert('Invalid email.');
                return;
            }
            if (password!==confirmPassword){
                alert('Password and confirm password are not the same.');
                return;
            }
            if (password.length<8) {
                alert('Password should contain at least 8 characters');
                return;
            }

            // Create an object with the form data
            const formData = {
                name,
                username,
                email,
                password,
            };
            console.log(JSON.stringify(formData))
            // Send a POST request to the server
            try {
                const response = await fetch('/auth/sign-up', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });

                if (response.ok) {
                    // Registration was successful
                    console.log('Registration successful');
                    const data = await response.json();
                    const jwtToken = data.token;
                    localStorage.setItem('token', jwtToken);
                    window.location.assign('/auth/recommendation')
                    // You can handle the success case here (e.g., show a success message or redirect the user).
                } else {
                    // Registration failed
                    console.error('Registration failed');

                    // Handle the failure case (e.g., show an error message to the user).
                }
            } catch (error) {
                console.error('Error:', error);
                // Handle network or other errors here.
            }
        });
    }
});
