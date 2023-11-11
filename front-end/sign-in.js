document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior

            // Gather form data
            const username = document.querySelector('input[name="username"]').value;
            const password = document.querySelector('input[name="password"]').value;

            // Create an object with the form data
            const formData = {
                username,
                password,
            };

            // Send a POST request to the server
            try {
                const response = await fetch('/sign-in', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(formData),
                });


                if (response.ok) {
                    // Registration was successful
                    const data = await response.json();
                    //
                    // // Извлечение JWT-токена из объекта данных
                    const jwtToken = data.token;
                    localStorage.setItem('token',jwtToken);
                    // Пример использования токена
                    console.log('Successful');
                    window.location.assign('/')
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
