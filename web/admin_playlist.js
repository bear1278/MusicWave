document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');


    const playlistListContainer = document.getElementById('playlist-list');

    fetch('/admin/playlist', {
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
            const playlists = data.playlists;

            // Используем функцию для создания элемента плейлиста
            function createPlaylistElement(playlist) {
                const playlistElement = document.createElement('div');
                playlistElement.classList.add('track-reason');
                console.log(playlist)
                playlistElement.innerHTML = `
                    <div class="track">
                        <div class="title">
                            <img class="track-icon" src="${playlist.cover}" alt="icon">
                            <div class="names">
                                <p>Name: </p>
                                <p>Author: </p>
                                <p>Type: </p>
                                <p>Release Date: </p>
                                <p>Modified Date: </p>
                            </div>
                            <div class="track-names">
                                <p><a class="link-names" href="http://localhost:8000/api/playlists/page/${playlist.id}">${playlist.name}</a></p>
                                <p>${playlist.Author.username}</p>
                                <p>${playlist.type}</p>
                                <p>${playlist.modifiedDate}</p>
                                <p>${playlist.releaseDate}</p>
                            </div>
                        </div>
                        <div class="time-buttons">
                            <div class="time-info">
                                <img src="/static/images/time.svg" class="timer" alt="duration:">
                                <p class="time">${formatDuration(playlist.duration)}</p>
                            </div>
                            <button class="button-delete playlist" title="Удалить плейлист"></button>
                        </div>
                    </div>
                    <form class="reason-input">
                        <label>
                            <input type="text" name="reason" placeholder="Reason" class="search-input">
                        </label>
                        <button type="submit" class="button-reason">Send</button>
                    </form>
                `;
                return playlistElement;
            }



            if (playlists === null) {
                const noPlaylistsMessage = document.createElement('div');
                noPlaylistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't playlists to suggest you.`;
                noPlaylistsMessage.classList.add('text-alert');

                noPlaylistsMessage.appendChild(messageParagraph);
                playlistListContainer.appendChild(noPlaylistsMessage);
            }else{
                // Создаем элементы для каждого плейлиста и добавляем их в контейнер
                playlists.forEach(playlist => {
                    const playlistElement = createPlaylistElement(playlist);
                    playlistListContainer.appendChild(playlistElement);

                    const button = playlistElement.querySelector('.button-delete');
                    const formSend = playlistElement.querySelector('.reason-input');


                    var jsonData
                    // Проверяем, есть ли значение
                    if (playlist) {
                        // Парсим JSON из значения data-value
                        button.addEventListener("click", function () {
                            var form = playlistElement.querySelector('.reason-input')
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
                            handlePostButtonLogic(button, jsonObject,`/admin/playlist/${playlist.id}`);
                        });

                    } else {
                        console.error('Data-value is missing for the button.');
                    }

                });
            }

            function handlePostButtonLogic(button, jsonData,url) {
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
            // Функция для форматирования длительности из миллисекунд в формат mm:ss
            function formatDuration(duration) {
                if (duration === 0 || duration === undefined) {
                    return `0:00`;
                }
                const minutes = Math.floor(duration / 60000);
                const seconds = ((duration % 60000) / 1000).toFixed(0);
                return `${minutes}:${(seconds < 10 ? '0' : '')}${seconds}`;
            }

        });
})
