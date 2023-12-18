document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior

            // Gather form data
            const username = document.querySelector('input[name="username"]').value;
            const email = document.querySelector('input[name="email"]').value;

            // Regular expressions for validation
            const usernameRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const emailRegex = /^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$/;

            // Perform validation
            if (!usernameRegex.test(username)) {
                alert('Invalid username. Use only Latin characters and special symbols.');
                return;
            }

            if (!emailRegex.test(email)) {
                alert('Invalid email.');
                return;
            }


            // Create an object with the form data
            const formData = {
                username,
                email,
            };

            // Send a POST request to the server
            try {
                const response = await fetch('/auth/reset-pass-email', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });


                if (response.ok) {
                    alert('We send link on email to reset password')
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
