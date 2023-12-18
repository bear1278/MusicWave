document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const UserListContainer = document.getElementById('user');

    function showError(error) {
        if (error && error.message) {
            alert(error.message);
        } else {
            alert("Произошла ошибка");
        }
    }

    document.getElementById("change-username").addEventListener("click", function () {
        var form = document.getElementById("change-username-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });
    document.getElementById("change-email").addEventListener("click", function () {
        var form = document.getElementById("change-email-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });
    document.getElementById("change-password").addEventListener("click", function () {
        var form = document.getElementById("change-password-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });
    document.getElementById("change-picture").addEventListener("click", function () {
        var form = document.getElementById("change-picture-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });




    document.getElementById('change-username-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);

        // Convert FormData to JSON
        var jsonObject = {};
        formData.forEach(function (value, key) {
            jsonObject[key] = value;
        });

        fetch('/profile/change-name', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization-1': 'Bearer ' + token,
            },
            body: JSON.stringify(jsonObject)
        })
            .then(response => {
                if (!response.ok) {
                    // Если код ответа не 200, отображаем сообщение в alert
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        alert(`Error: ${errorMessage}`);
                        throw new Error(errorMessage);
                    });
                }
                // Если код ответа 200, возвращаем данные
                return response.json();
            })
            .then(data => {
                // Обработка данных, если необходимо
                console.log(data);
                alert('Success.')
                window.location.reload(true);
            })
            .catch(error => {
                // Если не удалось распарсить JSON или нет поля "message"
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });

        location.reload(true)
    });


    document.getElementById('change-email-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);

        // Convert FormData to JSON
        var jsonObject = {};
        formData.forEach(function (value, key) {
            jsonObject[key] = value;
        });

        fetch('/profile/change-email', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token, // Исправление: 'Authorization-1' -> 'Authorization'
            },
            body: JSON.stringify(jsonObject),
        })
            .then(response => {
                if (!response.ok) {
                    // Если код ответа не 200, отображаем сообщение в alert
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        alert(`Error: ${errorMessage}`);
                        throw new Error(errorMessage);
                    });
                }
                // Если код ответа 200, возвращаем данные
                return response.json();
            })
            .then(data => {
                // Обработка данных, если необходимо
                console.log(data);
                alert('Success.')
                window.location.reload(true);
            })
            .catch(error => {
                // Если не удалось распарсить JSON или нет поля "message"
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });


    });


    document.getElementById('change-password-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);

        // Convert FormData to JSON
        var jsonObject = {};
        formData.forEach(function (value, key) {
            jsonObject[key] = value;
        });

        fetch('/profile/change-pass', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
                'Authorization-1': 'Bearer ' + token,
            },
            body: JSON.stringify(jsonObject)
        })
            .then(response => {
                if (!response.ok) {
                    // Если код ответа не 200, отображаем сообщение в alert
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        throw new Error(errorMessage);
                    });
                }
                // Если код ответа 200, возвращаем данные
                return response.json();
            })
            .then(data => {
                // Обработка данных, если необходимо
                console.log(data);
                alert('Success.')
                window.location.reload(true);
            })
            .catch(error => {
                // Если не удалось распарсить JSON или нет поля "message"
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });

    });



    document.getElementById('change-picture-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);



        fetch('/profile/change-picture', {
            method: 'PATCH',
            headers: {
                'Authorization-1': 'Bearer ' + token,
            },
            body: formData
        })
            .then(response => {
                if(response.ok){
                    alert('Success.');
                    location.reload(true)
                }
                else{
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        throw new Error(errorMessage);
                    });
                }
            })
            .catch(error => {
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });


    });


    fetch('/profile/user', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
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
            const user = data.user;

            // Используем функцию для создания элемента трека
            function createUserElement(user) {
                const userElement = document.createElement('div');
                userElement.classList.add('user-head');

                userElement.innerHTML = `
                <img src="${user.picture}" alt="user-picture" class="user-picture">
                <div class="user-head-text">
                    <h3 class="profile-header">Profile</h3>
                    <h1 class="username">${user.username}</h1>
                    <h4 class="number-of-artist">${user.email}</h4> <!--count of artists-->
                </div>
                `;
                return userElement;
            }

            if (user === null) {
                const noUserMessage = document.createElement('div');
                noUserMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `Error.`;
                noUserMessage.classList.add('text-alert');

                noUserMessage.appendChild(messageParagraph);
                UserListContainer.appendChild(noUserMessage);
            } else {
                // Создаем элементы для каждого альбома и добавляем их в контейнер

                const userElement = createUserElement(user);
                UserListContainer.appendChild(userElement);



            }
        });


    const trackListContainer = document.getElementById('track-list');
    const artistListContainer = document.getElementById('artist-list');

    fetch('/profile/spotify', {
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
            const artistsArray = data.artists;


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


            // Используем функцию для создания элемента артиста
            function createArtistElement(artist) {
                const artistElement = document.createElement('div');
                artistElement.classList.add('artist');

                artistElement.innerHTML = `
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
                        <a class="button-spotify" href="${artist.spotifyURL}" target="_blank"></a>
                        <button class="button-options" title="Добавить исполнителя в библиотеку" data-value='${JSON.stringify({
                    id: artist.id, name: artist.name,
                    picture: artist.picture, spotifyURL: artist.spotifyURL, popularity: artist.popularity
                })}'></button>
                    </div>
                `;
                return artistElement;
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


            if (artistsArray === null) {
                const noArtistsMessage = document.createElement('div');
                noArtistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't artists to suggest you.`;
                noArtistsMessage.classList.add('text-alert');

                noArtistsMessage.appendChild(messageParagraph);
                artistListContainer.appendChild(noArtistsMessage);
            }else{
                // Создаем элементы для каждого артиста и добавляем их в контейнер
                artistsArray.forEach(artist => {
                    const artistElement = createArtistElement(artist);
                    artistListContainer.appendChild(artistElement);

                    // Получаем кнопку из созданного элемента
                    const button = artistElement.querySelector('.button-options');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        // Парсим JSON из значения data-value
                        jsonData = JSON.parse(dataValue);
                        handlePostButtonLogic(button, jsonData,`/api/artists`);
                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }
                });
            }


            function handlePostButtonLogic(button, jsonData,url) {
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
                        then(postResponse => {
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
            // Функция для форматирования длительности из миллисекунд в формат mm:ss
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
        });




    const trackProfileListContainer = document.getElementById('track-list-profile');
    const artistProfileListContainer = document.getElementById('artist-list-profile');
    const durationListContainer = document.getElementById('total-duration');
    const genreListContainer = document.getElementById('genres');

    fetch('/profile/info', {
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
            const artistsArray = data.artists;
            const duration =data.duration;
            const genres=data.genres

            const genresName = genres.map(genre => `${genre.name}`).join(', ');

            const genreElement = document.createElement('div');
            genreElement.innerHTML=`<p class="link-names-track">${genresName}</p>`
            genreListContainer.appendChild(genreElement);

            const durationElement = document.createElement('div');
            durationElement.innerHTML=`<p class="link-names-track">${formatHoursDuration(duration)}</p>`
            durationListContainer.appendChild(durationElement);

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


            // Используем функцию для создания элемента артиста
            function createArtistElement(artist) {
                const artistElement = document.createElement('div');
                artistElement.classList.add('artist');

                artistElement.innerHTML = `
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
                        <a class="button-spotify" href="${artist.spotifyURL}" target="_blank"></a>
                        <button class="button-options" title="Добавить исполнителя в библиотеку" data-value='${JSON.stringify({
                    id: artist.id, name: artist.name,
                    picture: artist.picture, spotifyURL: artist.spotifyURL, popularity: artist.popularity
                })}'></button>
                    </div>
                `;
                return artistElement;
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


                    tracks.forEach(track => {
                        const trackElement = createTrackElement(track);
                        trackProfileListContainer.appendChild(trackElement);
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


            if (artistsArray === null) {
                const noArtistsMessage = document.createElement('div');
                noArtistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't artists to suggest you.`;
                noArtistsMessage.classList.add('text-alert');

                noArtistsMessage.appendChild(messageParagraph);
                artistProfileListContainer.appendChild(noArtistsMessage);
            }else{
                // Создаем элементы для каждого артиста и добавляем их в контейнер
                artistsArray.forEach(artist => {
                    const artistElement = createArtistElement(artist);
                    artistProfileListContainer.appendChild(artistElement);

                    // Получаем кнопку из созданного элемента
                    const button = artistElement.querySelector('.button-options');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        // Парсим JSON из значения data-value
                        jsonData = JSON.parse(dataValue);
                        handlePostButtonLogic(button, jsonData,`/api/artists`);
                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }
                });
            }


            function handlePostButtonLogic(button, jsonData,url) {
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
                });
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

            function formatHoursDuration(duration) {
                if (duration === 0 || duration === undefined) {
                    return `0:00:00`;
                }
                const hours = Math.floor(duration / 3600000);
                const minutes = Math.floor((duration % 3600000) / 60000);
                const seconds = ((duration % 60000) / 1000).toFixed(0);
                return `${(hours < 10 ? '0' : '')}${hours}:${(minutes < 10 ? '0' : '')}${minutes}:${(seconds < 10 ? '0' : '')}${seconds}`;
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
        });


})