(rift 
  (reference hello_rift) 
  (assignment 
    (reference a) 
    (integer 10)) 
  (assignment 
    (reference b) 
    (integer 20)) 
  (assignment 
    (reference sum) 
    (function 
      (arguments 
        (reference a) 
        (reference b)) 
      (operation 
        (reference a) 
        (binary-operator +) 
        (reference b)))) 
  (assignment 
    (reference c) 
    (function-apply
      (reference sum) 
      (tuple 
        (reference a) 
        (reference b)))))