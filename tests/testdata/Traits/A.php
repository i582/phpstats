<?php

trait Hello {
  public function sayHello() {
    echo 'Hello ';
  }
}

trait World {
  public function sayWorld() {
    echo 'World';
  }
}

class MyHelloWorld {
  use Hello, World;
  public function sayExclamationMark() {
    echo '!';
  }
}
