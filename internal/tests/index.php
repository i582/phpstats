<?php

class Foo
{
    const MAX_PASSWORD_LENGTH = 7; // not const SEVEN = 7 :)

    public function setPassword($password = 100)
    {
         if ($a == 100) {
            return 5;
         }

         $is_even = $number % 2 === 0;


         // don't do this
         if (mb_strlen($password) > 7) {
              throw new InvalidArgumentException("password");
         }
    }
}
