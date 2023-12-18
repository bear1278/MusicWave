document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const artistListContainer = document.getElementById('artist-list');


    fetch('/admin/artist', {
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

            const artistsArray = data.artists;


            function createArtistElement(artist) {
                const artistElement = document.createElement('div');
                artistElement.classList.add('artist-reason');

                artistElement.innerHTML = `
                    <div class="artist">
                        <div class="title">
                            <img class="track-icon" src="${artist.picture}" alt="icon">
                            <div class="artist-title">
                                <div class="names">
                                    <p>Artist: </p>
                                </div>
                                <div class="artist-names">
                                    <p><a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a></p>
                                </div>
                            </div>
                        </div>
                        <div class="buttons">
                            <div class="time-user">
                                <div>
                                    <p >Added Date:</p>
                                </div>
                                <div class="artist-names-date">
                                    <p class="link-names">${artist.addedDate}</p>
                                </div>
                            </div>
                            <a class="button-spotify" href="${artist.spotifyURL}" target="_blank"></a>
                            <button class="button-delete" title="Удалить исполнителя" ></button>
                        </div>
                     </div>
                    <form class="reason-input">
                        <label>
                            <input type="text" name="reason" placeholder="Reason" class="search-input">
                        </label>
                        <button type="submit" class="button-reason">Send</button>
                    </form>
                `;
                return artistElement;
            }


            if (artistsArray === null) {
                const noArtistsMessage = document.createElement('div');
                noArtistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't artists to suggest you.`;
                noArtistsMessage.classList.add('text-alert');

                noArtistsMessage.appendChild(messageParagraph);
                artistListContainer.appendChild(noArtistsMessage);
            } else {
                // Создаем элементы для каждого артиста и добавляем их в контейнер
                artistsArray.forEach(artist => {
                    const artistElement = createArtistElement(artist);
                    artistListContainer.appendChild(artistElement);

                    // Получаем кнопку из созданного элемента
                    const button = artistElement.querySelector('.button-delete');
                    const formSend = artistElement.querySelector('.reason-input');

                    var jsonData
                    // Проверяем, есть ли значение
                    if (artist) {
                        button.addEventListener("click", function () {
                            var form = artistElement.querySelector('.reason-input')
                            if (form.style.display === "none" || form.style.display === "") {
                                form.style.display = "flex";
                            } else {
                                form.style.display = "none";
                            }
                        });

                        formSend.addEventListener('submit', function (event) {
                            event.preventDefault();

                            var formData = new FormData(event.target);

                            // Convert FormData to JSON
                            var jsonObject = {};
                            formData.forEach(function (value, key) {
                                jsonObject[key] = value;
                            });
                            handlePostButtonLogic(button, jsonObject,`/admin/artist/${artist.id}`);
                        });

                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }
                });
            }


            function handlePostButtonLogic(button, jsonData, url) {
                // Добавляем обработчик события к кнопке для выполнения POST-запроса


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
        })
})
