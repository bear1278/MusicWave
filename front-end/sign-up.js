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
                const response = await fetch('/sign-up', {
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
                    //
                    // // Извлечение JWT-токена из объекта данных
                    const jwtToken = data.token;
                    localStorage.setItem('token', jwtToken);
                    window.location.assign('/recommendation')
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
