// Общие функции JavaScript для всех страниц

// Функция для форматирования чисел с двумя десятичными знаками
function formatNumber(num) {
    return num.toFixed(2);
}

// Функция для обновления времени на странице
function updateTime() {
    const timeElements = document.querySelectorAll('.current-time');
    const now = new Date();
    const formattedTime = now.toLocaleTimeString();
    
    timeElements.forEach(element => {
        element.textContent = formattedTime;
    });
}

// Обновление времени каждую секунду
setInterval(updateTime, 1000);

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    console.log('Страница загружена');
    updateTime();
});