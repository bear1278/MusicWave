document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const trackListContainer = document.getElementById('album-list');
    const albumListContainer =document.getElementById('user');
    var currentUrl = window.location.href;

// Разбиваем URL по слешам
    var urlParts = currentUrl.split('/');

// Получаем последнюю часть URL (последний элемент массива)
    var albumId = urlParts[urlParts.length - 1];
    console.log(albumId)


    fetch(`/api/albums/${albumId}`, {
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
            const album = data.album;
            console.log(album)


            // Используем функцию для создания элемента трека
            function createTrackElement(album) {
                const trackElement = document.createElement('div');
                trackElement.classList.add('track-head');

                const artistNames = album.artists.map(artist => `<a class="link-names-track" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a>`).join(', ');

                trackElement.innerHTML = `
                    <div class="track-head-text">
                        <img src="${album.cover}" alt="user-picture" class="user-picture">
                        <div class="names-head">
                            <p class="names-title">Name: </p>
                            <p class="names-title">Artists: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names-track" href="http://localhost:8000/api/tracks/page/${album.id}">${album.name}</a></p>
                            <p>${artistNames}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(album.duration)}</p>
                        </div>
                        <div class="buttons">
                            <a class="button-spotify" href="${album.spotifyURL}" target="_blank"></a>
                            <button class="button-options" title="Добавить альбом в библиотеку"  data-value='${JSON.stringify({id: album.id,
                    name: album.name, cover: album.cover, releaseDate: album.releaseDate, albumType: album.albumType,
                    totalTracks: album.totalTracks, spotifyURL: album.spotifyURL, popularity: album.popularity,
                    artists: album.artists})}'></button>
                            </div>
                        </div>
                    </div>
                    `;

                return trackElement;
            }

            if (album === null) {
                const noAlbumsMessage = document.createElement('div');
                noAlbumsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't albums to suggest you.`;
                noAlbumsMessage.classList.add('text-alert');

                noAlbumsMessage.appendChild(messageParagraph);
                albumListContainer.appendChild(noAlbumsMessage);
            }else{
                // Создаем элементы для каждого альбома и добавляем их в контейнер

                    const albumElement = createTrackElement(album);
                    albumListContainer.appendChild(albumElement);

                    const button = albumElement.querySelector('.button-options');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        // Парсим JSON из значения data-value
                        jsonData = JSON.parse(dataValue);
                        handlePostButtonLogic(button, jsonData,`/api/albums/`);
                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }



            }


        })
    fetch(`/api/albums/${albumId}/tracks`, {
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
                            <a class="button-spotify" href="${track.spotifyURL}" target="_blank"></a>
                            <button class="button-track" title="Добавить трек в плейлист" data-value='${JSON.stringify({type: "track", id: track.id})}'></button>
                            <div class="context-menu" id="contextMenu">
                                
                            </div>
                        </div>
                    </div>
                `;
                return trackElement;
            }

            fetch('/api/playlists/my', {
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
                        throw new Error(`HTTP error! Status: ${response.status} - ${response.statusText}`);
                    }
                    return response.json();
                })
                .then(playlists => {

                    console.log(tracks)
                    tracks.forEach(track => {
                        const trackElement = createTrackElement(track);
                        trackListContainer.appendChild(trackElement);
                        const button = trackElement.querySelector('.button-track');
                        const contextMenu = trackElement.querySelector('.context-menu');

                        playlists.playlists.forEach(playlist => {
                            const playlistButton = document.createElement('button');
                            playlistButton.textContent = playlist.name;
                            playlistButton.classList.add('playlist-button');
                            playlistButton.addEventListener('click', function () {
                                // Шаг 3: Отправляем POST-запрос при нажатии на кнопку плейлиста
                                handleAddToPlaylist(playlist.id, track);
                                // Удаляем кнопку плейлиста из контекстного меню
                                contextMenu.removeChild(playlistButton);
                            });

                            contextMenu.appendChild(playlistButton);


                            let isContextMenuVisible = false;

                            // Прикрепляем обработчик события на кнопку трека
                            button.addEventListener('click', function (event) {
                                event.preventDefault();

                                // Переключаем состояние видимости контекстного меню
                                isContextMenuVisible = !isContextMenuVisible;

                                // Позиционируем контекстное меню около кнопки
                                const buttonRect = button.getBoundingClientRect();
                                contextMenu.style.left = buttonRect.left + '-px';
                                contextMenu.style.top = buttonRect.bottom + '-px';
                                // Показываем или скрываем контекстное меню в зависимости от состояния
                                contextMenu.style.display = isContextMenuVisible ? 'block' : 'none';

                                // Предотвращаем скрытие контекстного меню при клике внутри него
                                contextMenu.addEventListener('click', function (event) {
                                    event.stopPropagation();
                                });
                            });


                        });

                    });
                })
        });
    function handleAddToPlaylist(playlistId, track) {
        fetch(`/api/playlists/${playlistId}/tracks`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization-1': 'Bearer ' + token,
                'Authorization-2': 'Bearer ' + accessToken,
                'Authorization-3': 'Bearer ' + refreshToken,
                'Authorization-4': 'Bearer ' + expiry,
            },
            body: JSON.stringify({ id: track.id, name: track.name, duration: track.duration,
                cover: track.cover,album: track.album,spotifyURL: track.spotifyURL, popularity: track.popularity}),
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
})

function formatDuration(duration) {
    if (duration === 0 || duration === undefined) {
        return `0:00`;
    }
    const minutes = Math.floor(duration / 60000);
    const seconds = ((duration % 60000) / 1000).toFixed(0);
    return `${minutes}:${(seconds < 10 ? '0' : '')}${seconds}`;
}

function handlePostButtonLogic(button, jsonData,url) {
    // Добавляем обработчик события к кнопке для выполнения POST-запроса
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');
    button.addEventListener('click', function (event) {
        event.preventDefault();
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
    });

}

