(define (pascal row col)
    (if (on-edge? row col)
      1
      (+ (pascal (- row 1) (- col 1))
         (pascal (- row 1) col))))

(define (on-edge? row col)
  (or (= col 1)
      (= row col)))

(display (pascal 10 3)) ; returns 36
(newline)

(define (factorial n)
  (if (= n 1)
    1
    (* (factorial (- n 1)) n)))

(display (factorial 5)) ; returns 120
(newline)
