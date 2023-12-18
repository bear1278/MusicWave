document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const userListContainer = document.getElementById('user-list');


    fetch('/admin/user', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
            'Authorization-2': 'Bearer ' + accessToken,
            'Authorization-3': 'Bearer ' + refreshToken,
            'Authorization-4': 'Bearer ' + expiry,
        },
    })
        .then(response => {
            if (!response.ok) {
                if (response.status === 401) {
                    alert('You need to enter to your account.');
                    window.location.assign('/auth/sign-in');
                }
                throw new Error(`HTTP error! Status: ${response.status} - ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            const users = data.users;



            // Используем функцию для создания элемента альбома
            function createUserElement(user) {
                const userElement = document.createElement('div');
                userElement.classList.add('track-reason');


                userElement.innerHTML = `
                    <div class="track">
                        <div class="title">
                            <img class="track-icon" src="${user.picture}" alt="icon">
                            <div class="names">
                                <p>Name: </p>
                                <p>Username: </p>
                                <p>Email: </p>
                                <p>Create Date: </p>
                                <p>Modified Date: </p>
                            </div>
                            <div class="track-names">
                                <p>${user.name}</p>
                                <p>${user.username}</p>
                                <p>${user.email}</p>
                                <p>${user.createDate}</p>
                                <p>${user.modifiedDate}</p>
                            </div>
                        </div>
                        <div class="user-buttons">
                            <div class="buttons-user">
                                <button class="button-delete-user" title="Удалить пользователя"></button>
                            </div>
                        </div>
                    </div>
                    <form class="reason-input">
                        <label>
                            <input type="text" name="reason" placeholder="Reason" class="search-input">
                        </label>
                        <button type="submit"  class="button-reason">Send</button>
                    </form>
                `;
                return userElement;
            }


            if (users === null) {
                const noAlbumsMessage = document.createElement('div');
                noAlbumsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't albums to suggest you.`;
                noAlbumsMessage.classList.add('text-alert');

                noAlbumsMessage.appendChild(messageParagraph);
                userListContainer.appendChild(noAlbumsMessage);
            }else{
                // Создаем элементы для каждого альбома и добавляем их в контейнер
                users.forEach(user => {
                    console.log(user)
                    const userElement = createUserElement(user);
                    userListContainer.appendChild(userElement);

                    const button = userElement.querySelector('.button-delete-user');
                    const SendForm =userElement.querySelector('.reason-input')

                    var jsonData
                    // Проверяем, есть ли значение
                    if (user) {
                        // Парсим JSON из значения data-value
                        button.addEventListener("click", function () {
                            var form = userElement.querySelector('.reason-input')
                            if (form.style.display === "none" || form.style.display === "") {
                                form.style.display = "flex";
                            } else {
                                form.style.display = "none";
                            }
                        });

                        SendForm.addEventListener('submit', function (event) {
                            event.preventDefault();

                            var formData = new FormData(event.target);

                            // Convert FormData to JSON
                            var jsonObject = {};
                            formData.forEach(function (value, key) {
                                jsonObject[key] = value;
                            });
                            handlePostButtonLogic(button, jsonObject,`/admin/user/${user.id}`);
                        });


                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }


                });
            }


            function handlePostButtonLogic(button, jsonData,url) {
                // Добавляем обработчик события к кнопке для выполнения POST-запроса

                // Проверяем наличие data-value
                if (jsonData) {
                    // Выполняем POST-запрос
                    fetch(url, {
                        method: 'DELETE',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization-1': 'Bearer ' + token,
                            'Authorization-2': 'Bearer ' + accessToken,
                            'Authorization-3': 'Bearer ' + refreshToken,
                            'Authorization-4': 'Bearer ' + expiry,
                        },
                        body: JSON.stringify(jsonData),
                    })
                        .then(postResponse => {
                            if (postResponse.ok) {
                                console.log('POST request successful');
                                button.style.background = 'url("../../static/images/check.svg")';
                                button.disabled = true;
                                window.location.reload(true)
                            } else {
                                throw new Error(`HTTP error! Status: ${postResponse.status} - ${postResponse.statusText}`);
                            }
                        })
                        .catch(error => {
                            console.error('Error sending POST request:', error);
                            // Добавьте свою логику обработки ошибок POST-запроса
                        });
                } else {
                    console.error('Data-value is missing for the button.');
                }
            }

        });
})
