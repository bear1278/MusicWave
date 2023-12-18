document.addEventListener('DOMContentLoaded', function() {

    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    var nextPage
    var previous

    const ListContainer = document.getElementById('objects');
    const ButtonListContainer = document.getElementById('buttons');

    document.getElementById('button-search').addEventListener('click', performSearch);

    function clearListContainer() {
        ListContainer.innerHTML = '';
        ButtonListContainer.innerHTML = '';
    }

    function appendButtons(url1,url2,searchType) {
        const buttonContainer = document.createElement('div');
        buttonContainer.classList.add('button-container');

        const button1 = document.createElement('button');
        button1.classList.add('button-page-back')
        button1.addEventListener('click', function() {
            if (url1!==''){
                clearListContainer()
                sendRequest(url1, searchType); // Adjust the URL and searchType as needed
            }
        });

        console.log(url2)
        const button2 = document.createElement('button');
        button2.classList.add('button-page-next')
        button2.addEventListener('click', function() {
            if (url2!==''){
                clearListContainer()
                sendRequest(url2, searchType); // Adjust the URL and searchType as needed
            }
        });

        buttonContainer.appendChild(button1);
        buttonContainer.appendChild(button2);

        ButtonListContainer.appendChild(buttonContainer);
    }

    function performSearch() {
        var searchInput = document.querySelector('.search-input');
        var selectedRadio = document.querySelector('input[name="object"]:checked');

        if (searchInput.value && selectedRadio) {
            var searchType = selectedRadio.value;
            var searchText = encodeURIComponent(searchInput.value);
            var url
            // Construct the URL
            if (searchType!=='playlists'){
                url = '/search/' + searchType + '/' + searchText + '/0';
            }else{
                url = '/search/' + searchType + '/' + searchText;
            }
            clearListContainer()
            // Perform AJAX request
            sendRequest(url, searchType)


        } else {
            alert('Please enter search text and select a search type.');
        }
    }

    function sendRequest(url,searchType) {
         fetch(url, {
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
                try {
                    const objects = data.data;
                    console.log(objects)
                    nextPage = data.nextPage;
                    previous = data.previous;

                    function createTrackElement(object) {
                        const trackElement = document.createElement('div');
                        trackElement.classList.add('track');

                        const artistNames = object.album.artists.map(artist => `<a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a>`).join(', ');

                        trackElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${object.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Artists: </p>
                            <p>Album: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/tracks/page/${object.id}">${object.name}</a></p>
                            <p>${artistNames}</a></p>
                            <p><a class="link-names" href="http://localhost:8000/api/albums/page/${object.album.id}">${object.album.name}</a></p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(object.duration)}</p>
                        </div>
                        <div class="buttons">
                            <a class="button-spotify" href="${object.spotifyURL}" target="_blank"></a>
                            <button class="button-track" title="Добавить трек в плейлист" data-value='${JSON.stringify({type: "track", id: object.id})}'></button>
                            <div class="context-menu" id="contextMenu">
                                
                            </div>
                        </div>
                    </div>
                `;
                        return trackElement;
                    }


                    // Используем функцию для создания элемента альбома
                    function createAlbumElement(object) {
                        const albumElement = document.createElement('div');
                        albumElement.classList.add('track');



                        for (let key in object) {
                            if (object.hasOwnProperty(key) && typeof object[key] === 'string') {
                                // Заменяем символы ' на \'
                                object[key] = object[key].replace(/'/g, "\\'");
                            }
                        }

                        var jsonData=JSON.stringify({
                            id: object.id,
                            name: object.name,
                            cover: object.cover,
                            releaseDate: object.releaseDate,
                            albumType: object.albumType,
                            totalTracks: object.totalTracks,
                            spotifyURL: object.spotifyURL,
                            popularity: object.popularity,
                            artists: object.artists
                        })

                        const artistNames = object.artists.map(artist => `<a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a>`).join(', ');

                        albumElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${object.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Artists: </p>
                            <p>Type: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/albums/page/${object.id}">${object.name}</a></p>
                            <p>${artistNames}</p>
                            <p>${object.albumType}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        
                        <div class="buttons">
                            <a class="button-spotify" href="${object.spotifyURL}" target="_blank"></a>
                            <button class="button-options" title="Добавить альбом в библиотеку" data-value='${jsonData}'></button>
                        </div>
                    </div>
                `;
                        return albumElement;
                    }

                    // Используем функцию для создания элемента артиста
                    function createArtistElement(object) {
                        const artistElement = document.createElement('div');
                        artistElement.classList.add('artist');

                        artistElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${object.picture}" alt="icon">
                        <div class="artist-title">
                            <div class="names">
                                <p>Artist: </p>
                            </div>
                            <div class="artist-names">
                                <p><a class="link-names" href="http://localhost:8000/api/artists/page/${object.id}">${object.name}</a></p>
                            </div>
                        </div>
                    </div>
                    <div class="buttons">
                        <a class="button-spotify" href="${object.spotifyURL}" target="_blank"></a>
                        <button class="button-options" title="Добавить исполнителя в библиотеку" data-value='${JSON.stringify({
                            id: object.id, name: object.name,
                            picture: object.picture, spotifyURL: object.spotifyURL, popularity: object.popularity
                        })}'></button>
                    </div>
                `;
                        return artistElement;
                    }

                    // Используем функцию для создания элемента плейлиста
                    function createPlaylistElement(object) {
                        const playlistElement = document.createElement('div');
                        playlistElement.classList.add('track');
                        playlistElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${object.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Author: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/playlists/page/${object.id}">${object.name}</a></p>
                            <p>${object.Author.username}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(object.duration)}</p>
                        </div>
                        <button class="button-options playlist" title="Добавить плейлист в библиотеку" data-value='${object.id}'></button>
                    </div>
                `;
                        return playlistElement;
                    }

                    if (searchType === 'track') {
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


                                objects.forEach(object => {
                                    const trackElement = createTrackElement(object);
                                    ListContainer.appendChild(trackElement);
                                    const button = trackElement.querySelector('.button-track');
                                    const contextMenu = trackElement.querySelector('.context-menu');

                                    playlists.playlists.forEach(playlist => {
                                        const playlistButton = document.createElement('button');
                                        playlistButton.textContent = playlist.name;
                                        playlistButton.classList.add('playlist-button');
                                        playlistButton.addEventListener('click', function () {
                                            // Шаг 3: Отправляем POST-запрос при нажатии на кнопку плейлиста
                                            handleAddToPlaylist(playlist.id, object);
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
                    }

                    if (searchType === 'album') {
                        if (objects === null) {
                            const noAlbumsMessage = document.createElement('div');
                            noAlbumsMessage.classList.add('alert');

                            const messageParagraph = document.createElement('p');
                            messageParagraph.textContent = `Not Found.`;
                            noAlbumsMessage.classList.add('text-alert');

                            noAlbumsMessage.appendChild(messageParagraph);
                            ListContainer.appendChild(noAlbumsMessage);
                        } else {
                            // Создаем элементы для каждого альбома и добавляем их в контейнер
                            objects.forEach(object => {
                                const albumElement = createAlbumElement(object);
                                ListContainer.appendChild(albumElement);

                                const button = albumElement.querySelector('.button-options');

                                //var dataValue = button.getAttribute('data-value');
                                var jsonData
                                // Проверяем, есть ли значение
                                console.log(object)
                                if (object) {
                                    // Парсим JSON из значения data-value

                                    jsonData = JSON.parse(JSON.stringify(object))
                                    handlePostButtonLogic(button, jsonData, `/api/albums/`);
                                    // Выполняем общую логику, отправляем POST-запрос
                                } else {
                                    console.error('Data-value is missing for the button.');
                                }


                            });
                        }
                    }

                    if (searchType === 'artist') {
                        if (objects === null) {
                            const noArtistsMessage = document.createElement('div');
                            noArtistsMessage.classList.add('alert');

                            const messageParagraph = document.createElement('p');
                            messageParagraph.textContent = `Not Found.`;
                            noArtistsMessage.classList.add('text-alert');

                            noArtistsMessage.appendChild(messageParagraph);
                            ListContainer.appendChild(noArtistsMessage);
                        } else {
                            // Создаем элементы для каждого артиста и добавляем их в контейнер
                            objects.forEach(object => {
                                const artistElement = createArtistElement(object);
                                ListContainer.appendChild(artistElement);

                                // Получаем кнопку из созданного элемента
                                const button = artistElement.querySelector('.button-options');


                                var jsonData
                                // Проверяем, есть ли значение
                                if (object) {
                                    // Парсим JSON из значения data-value
                                    jsonData = JSON.parse(JSON.stringify(object))
                                    handlePostButtonLogic(button, jsonData, `/api/artists`);
                                    // Выполняем общую логику, отправляем POST-запрос
                                } else {
                                    console.error('Data-value is missing for the button.');
                                }
                            });
                        }
                    }

                    if (searchType === 'playlists') {
                        if (objects === null) {
                            const noPlaylistsMessage = document.createElement('div');
                            noPlaylistsMessage.classList.add('alert');

                            const messageParagraph = document.createElement('p');
                            messageParagraph.textContent = `Not Found.`;
                            noPlaylistsMessage.classList.add('text-alert');

                            noPlaylistsMessage.appendChild(messageParagraph);
                            ListContainer.appendChild(noPlaylistsMessage);
                        } else {
                            // Создаем элементы для каждого плейлиста и добавляем их в контейнер
                            objects.forEach(object => {
                                const playlistElement = createPlaylistElement(object);
                                ListContainer.appendChild(playlistElement);

                                const button = playlistElement.querySelector('.button-options');

                                var dataValue = button.getAttribute('data-value');
                                var jsonData
                                // Проверяем, есть ли значение
                                if (dataValue) {
                                    // Парсим JSON из значения data-value
                                    jsonData = JSON.parse(dataValue);
                                    handlePostButtonLogic(button, jsonData, `/api/playlists/add/${dataValue}`);
                                    // Выполняем общую логику, отправляем POST-запрос
                                } else {
                                    console.error('Data-value is missing for the button.');
                                }

                            });
                        }
                    }

                    function handlePostButtonLogic(button, jsonData, url) {
                        // Добавляем обработчик события к кнопке для выполнения POST-запроса

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
                                    .then(postResponse => {
                                        if (postResponse.ok) {
                                            console.log('POST request successful');
                                            button.style.background = 'url("../../static/images/check.svg")';
                                            button.disabled = true;
                                        } else {
                                            return postResponse.json().then(data => {
                                                const errorMessage = data.message || 'Unknown error occurred';
                                                throw new Error(errorMessage);
                                            });
                                        }
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

                    if (searchType !== 'playlists') {
                        appendButtons(previous, nextPage, searchType)
                    }
                }catch (error) {
                    console.error('Error parsing JSON:', error);
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
    }
})