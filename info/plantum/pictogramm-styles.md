# Для создания диаграмм с помощью Python-библиотеки Diagrams
    (diagrams.mingrammer.com) - пиктограммы и стили на следующих сайтах:

## Официальная документация Diagrams
* **[AWS](https://diagrams.mingrammer.com/docs/nodes/aws)**
* **[GCP](https://diagrams.mingrammer.com/docs/nodes/gcp)**
* **[AZURE](https://diagrams.mingrammer.com/docs/nodes/azure)**
* **[ONPREM](https://diagrams.mingrammer.com/docs/nodes/onprem)**

### Font Awesome
* **[fontawesome](https://fontawesome.com/icons)**
  * Предлагает множество иконок, которые можно использовать в диаграммах

### Material Design Icons
* **[materialdesignicons.com](https://materialdesignicons.com/)**
  * Коллекция иконок от Google

### flaticon
* **[Flaticon](https://www.flaticon.com/)**
  * Большая библиотека бесплатных иконок

### The Noun Project
* **[The Noun Project](https://thenounproject.com/)**
  * Миллионы иконок, созданных глобальным сообществом

### Icons8
* **[Icons8](https://icons8.com/)**
  * Предлагает иконки в различных стилях

### GitHub репозиторий Diagrams
* **[GitHub репозиторий Diagrams](https://github.com/mingrammer/diagrams)**
  * Содержит все доступные ноды и иконки, используемые в библиотеке

### Devicon
* **[Devicon](https://devicon.dev/)**
  * Набор иконок, представляющих языки программирования, фреймворки и инструменты разработки

## Для использования собственных иконок в Diagrams,
    вы можете создать пользовательские ноды, указав путь к вашим изображениям:
```python
from diagrams import Diagram, Node

class CustomNode(Node):
    def __init__(self, label, **kwargs):
        super().__init__(label, **kwargs)

with Diagram("Custom Diagram"):
    CustomNode("My Custom Node", icon_path="/path/to/your/icon.png")
```
