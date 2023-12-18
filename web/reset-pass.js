document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior

            // Gather form data
            const password = document.querySelector('input[name="password"]').value;
            const confirmPassword = document.querySelector('input[name="confirmPassword"]').value;

            // Regular expressions for validation
            const passwordRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;
            const confirmPasswordRegex = /^[a-zA-Z0-9!@#$%^&*()_+{}\[\]:;<>,.?~\\/-]+$/;

            // Perform validation
            if (!passwordRegex.test(password)) {
                alert('Invalid password. Use only Latin characters and special symbols.');
                return;
            }

            if (!confirmPasswordRegex.test(confirmPassword)) {
                alert('Invalid confirm password. Use only Latin characters and special symbols.');
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
            const token = window.location.pathname.replace('/auth/reset-pass/', '');


            // Create an object with the form data
            const formData = {
                password,
                token,
            };
            // Send a POST request to the server
            try {
                const response = await fetch('/auth/reset-pass', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });


                if (response.ok) {
                    window.location.assign('/auth/sign-in')
                } else {
                    // Registration failed

                    const responseBody = await response.text();
                    console.log(responseBody)
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
