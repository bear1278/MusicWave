document.addEventListener('DOMContentLoaded', () => {
    const form = document.querySelector('form');

    if (form) { // Check if the form element exists
        form.addEventListener('submit', async (e) => {
            e.preventDefault(); // Prevent the default form submission behavior


            var selectedGenres = [];
            var checkboxes = document.querySelectorAll('input[name="genre"]:checked');
            checkboxes.forEach(function(checkbox) {
                selectedGenres.push(parseInt(checkbox.value));
            });

            var data = { genres: selectedGenres };

            var token = localStorage.getItem('token');
            // Send a POST request to the server
            try {
                const response = await fetch('/auth/recommendation', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization-1': 'Bearer ' + token,
                    },
                    body: JSON.stringify(data),
                });

                if (response.ok) {
                    // Registration was successful
                    console.log('Registration successful');
                    window.location.assign('/login')
                    // You can handle the success case here (e.g., show a success message or redirect the user).
                } else {
                    // Registration failed
                    console.error('Registration failed');
                    console.log(JSON.stringify(data))
                    // Handle the failure case (e.g., show an error message to the user).
                }
            } catch (error) {
                console.error('Error:', error);
                // Handle network or other errors here.
            }
        });
    }
});
