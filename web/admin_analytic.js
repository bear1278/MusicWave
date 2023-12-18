document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');





    fetch('/admin/genre', {
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
        .then(dataJson => {
            const genres = dataJson.genres


            var data = [];
            for (var i = 0; i < genres.length; i++) {
                data.push({ genre: genres[i].name, popularity: genres[i].popularity });
            }
            data.sort(function(a, b) {
                return b.popularity - a.popularity;
            });

// Получаем контекст рисования для элемента canvas
            var ctx = document.getElementById('myChart').getContext('2d');

// Создаем объект графика
            var myChart = new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: data.map(item => item.genre), // Метки по оси X (жанры)
                    datasets: [{
                        label: 'Популярность',
                        data: data.map(item => item.popularity), // Данные для графика (популярность)
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true
                        }
                    }
                }
            });



        });

    fetch('/admin/artist/popularity', {
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
        .then(dataJson => {
            const artists = dataJson.artists


            var data = [];
            for (var i = 0; i < artists.length; i++) {
                data.push({ artist: artists[i].name, popularity: artists[i].popularity });
            }
            data.sort(function(a, b) {
                return b.popularity - a.popularity;
            });

// Получаем контекст рисования для элемента canvas
            var ctxArt = document.getElementById('myChartArtist').getContext('2d');

// Создаем объект графика
            var myChart = new Chart(ctxArt, {
                type: 'bar',
                data: {
                    labels: data.map(item => item.artist), // Метки по оси X (жанры)
                    datasets: [{
                        label: 'Популярность',
                        data: data.map(item => item.popularity), // Данные для графика (популярность)
                        backgroundColor: 'rgba(75, 192, 192, 0.2)',
                        borderColor: 'rgba(75, 192, 192, 1)',
                        borderWidth: 1
                    }]
                },
                options: {
                    scales: {
                        y: {
                            beginAtZero: true
                        }
                    }
                }
            });

        });


    fetch('/admin/genre/diversity', {
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
        .then(dataJson => {
            const genres = dataJson.genres


            var ctx = document.getElementById('myPieChart').getContext('2d');
            ctx.canvas.width = 500;
            ctx.canvas.height = 500;
            var myPieChart = new Chart(ctx, {
                type: 'pie',
                data: {
                    labels: genres.map(genre => genre.name), // Метки (жанры)
                    datasets: [{
                        data: genres.map(genre => genre.diversity), // Данные для круговой диаграммы (разнообразие)
                        backgroundColor: generateRandomColors(genres.length),
                        borderColor: '#ffffff', // Цвет границы
                        borderWidth: 1
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false
                }
            });

            // Функция для генерации случайных цветов
            function generateRandomColors(numColors) {
                var colors = [];
                for (var i = 0; i < numColors; i++) {
                    colors.push(getRandomColor());
                }
                return colors;
            }

            // Функция для получения случайного цвета
            function getRandomColor() {
                var letters = '0123456789ABCDEF';
                var color = '#';
                for (var i = 0; i < 6; i++) {
                    color += letters[Math.floor(Math.random() * 16)];
                }
                return color;
            }
        });
})
