document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const trackListContainer = document.getElementById('album-list');
    const playlistListContainer =document.getElementById('user');
    const ChangeContainer =document.getElementById('change-playlist')
    var currentUrl = window.location.href;
    var addSpotifyPermission
    var buttonAddSpotify=``

// Разбиваем URL по слешам
    var urlParts = currentUrl.split('/');

// Получаем последнюю часть URL (последний элемент массива)
    var playlistId = urlParts[urlParts.length - 1];

    function clearListContainer() {
        ChangeContainer.innerHTML = '';
    }

    document.getElementById('create-playlist-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);

        fetch(`/api/playlists/${playlistId}`, {
            method: 'PUT',
            headers: {
                'Authorization-1': 'Bearer ' + token,
            },
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                console.log('Ответ от сервера:', data);
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
        window.location.reload(true)
    });




    document.getElementById("create-playlist").addEventListener("click", function () {
        var form = document.getElementById("create-playlist-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });




    fetch(`/api/playlists/${playlistId}`, {
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
            const playlist = data.playlist;
            const changePermission= data.changePermission

            addSpotifyPermission=data.addSpotifyPermission

            if (!changePermission){
                clearListContainer()
            }

            // Используем функцию для создания элемента трека
            function createPlaylistElement(playlist) {
                const playlistElement = document.createElement('div');
                playlistElement.classList.add('track-head');


                playlistElement.innerHTML = `
                    <div class="track-head-text">
                        <img src="${playlist.cover}" alt="user-picture" class="user-picture">
                        <div class="names-head">
                            <p class="names-title">Name: </p>
                            <p class="names-title">Author: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names-track" href="http://localhost:8000/api/tracks/page/${playlist.id}">${playlist.name}</a></p>
                            <p class="link-names-track">${playlist.Author.username}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(playlist.duration)}</p>
                        </div>
                        <div class="buttons">
                            <button class="button-options" title="Добавить плейлист в библиотеку"></button>
                            </div>
                        </div>
                    </div>
                    `;

                return playlistElement;
            }

            if (playlist === null) {
                const noAlbumsMessage = document.createElement('div');
                noAlbumsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't albums to suggest you.`;
                noAlbumsMessage.classList.add('text-alert');

                noAlbumsMessage.appendChild(messageParagraph);
                playlistListContainer.appendChild(noAlbumsMessage);
            }else{
                // Создаем элементы для каждого альбома и добавляем их в контейнер

                const playlistElement = createPlaylistElement(playlist);
                playlistListContainer.appendChild(playlistElement);

                const button = playlistElement.querySelector('.button-options');

                var jsonData
                // Проверяем, есть ли значение
                if (playlist) {
                    // Парсим JSON из значения data-value
                    button.addEventListener('click', () => {
                        jsonData = JSON.parse(JSON.stringify(playlist));

                        handlePostButtonLogic(button, jsonData, `/api/playlists/`);
                    });
                    // Выполняем общую логику, отправляем POST-запрос
                } else {
                    console.error('Data-value is missing for the button.');
                }



            }


        })
    fetch(`/api/playlists/${playlistId}/tracks`, {
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
            const tracks = data.tracks;
            console.log(tracks)
            // Используем функцию для создания элемента трека
            function createTrackElement(track) {
                const trackElement = document.createElement('div');
                trackElement.classList.add('track');

                const artistNames = track.album.artists.map(artist => `<a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a>`).join(', ');
                if (addSpotifyPermission){
                    buttonAddSpotify=`
                    <button class="button-spotify-add">add to Spotify</button>
                    `
                }
                trackElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${track.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Artists: </p>
                            <p>Album: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/tracks/page/${track.id}">${track.name}</a></p>
                            <p>${artistNames}</a></p>
                            <p><a class="link-names" href="http://localhost:8000/api/albums/page/${track.album.id}">${track.album.name}</a></p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(track.duration)}</p>
                        </div>
                        <div class="buttons">
                            <a class="button-spotify" href="${track.spotifyURL}" target="_blank"></a>`
                            +buttonAddSpotify+
                            `<button class="button-delete" title="Удалить трек" data-value='${JSON.stringify({type: "track", id: track.id})}'></button>
                        </div>
                    </div>
                `;
                return trackElement;
            }

            if (tracks === null) {
                const noAddedPlaylistsMessage = document.createElement('div');
                noAddedPlaylistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `You haven't got added playlists.`;
                noAddedPlaylistsMessage.classList.add('text-alert');

                noAddedPlaylistsMessage.appendChild(messageParagraph);
                trackListContainer.appendChild(noAddedPlaylistsMessage);
            }else {
                tracks.forEach(track => {
                    const trackElement = createTrackElement(track);
                    trackListContainer.appendChild(trackElement);

                    const button = trackElement.querySelector('.button-delete');

                    const buttonSpotify = trackElement.querySelector('.button-spotify-add');

                    var jsonData
                    // Проверяем, есть ли значение
                    if (track) {
                        // Парсим JSON из значения data-value
                        button.addEventListener('click', () => {

                            jsonData = JSON.parse(JSON.stringify(track));
                            handleDeleteButtonLogic(button, jsonData, `/api/playlists/${playlistId}/tracks/${track.id}`);
                            // Выполняем общую логику, отправляем POST-запрос
                        });
                        if (addSpotifyPermission) {
                            buttonSpotify.addEventListener('click', () => {
                                const userConfirmed = confirm('Вы хотите добавить песню в свой аккаунт спотифай?');
                                if (userConfirmed) {
                                    jsonData = JSON.parse(JSON.stringify(track));
                                    handlePostButtonLogic(buttonSpotify, jsonData, `/api/tracks/${track.id}`);
                                    // Выполняем общую логику, отправляем POST-запрос
                                }
                            });
                        }
                    } else {
                        console.error('Data-value is missing for the button.');
                    }

                });
            }
        });
    function handleDeleteButtonLogic(button, jsonData,url) {
        // Добавляем обработчик события к кнопке для выполнения POST-запроса


        // Проверяем наличие data-value

        // Выполняем POST-запрос
        fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization-1': 'Bearer ' + token,
            },
        })
            .then(response => {
                if (!response.ok) {
                    // Если код ответа не 200, отображаем сообщение в alert
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        throw new Error(errorMessage);
                    });
                }
                console.log('Track added to playlist successfully');
            })
            .catch(error => {
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });

    }

    function handlePostButtonLogic(button, jsonData,url) {
        // Добавляем обработчик события к кнопке для выполнения POST-запроса

            // Проверяем наличие data-value
            if (jsonData) {
                // Выполняем POST-запрос
                fetch(url, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization-1': 'Bearer ' + token,
                        'Authorization-2': 'Bearer ' + accessToken,
                        'Authorization-3': 'Bearer ' + refreshToken,
                        'Authorization-4': 'Bearer ' + expiry,
                    },
                    body: JSON.stringify(jsonData),
                })
                    .then(response => {
                        if (!response.ok) {
                            // Если код ответа не 200, отображаем сообщение в alert
                            return response.json().then(data => {
                                const errorMessage = data.message || 'Unknown error occurred';
                                throw new Error(errorMessage);
                            });
                        }
                        console.log('Track added to playlist successfully');
                    })
                    .catch(error => {
                        alert(`Error: ${error.message || 'Unknown error occurred'}`);
                        console.error(error);
                    });
            } else {
                console.error('Data-value is missing for the button.');
            }

    }






})

function formatDuration(duration) {
    if (duration === 0 || duration === undefined) {
        return `0:00`;
    }
    const minutes = Math.floor(duration / 60000);
    const seconds = ((duration % 60000) / 1000).toFixed(0);
    return `${minutes}:${(seconds < 10 ? '0' : '')}${seconds}`;
}

