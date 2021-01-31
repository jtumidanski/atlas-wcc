package producers

import "strconv"

func createKey(key int) []byte {
   var empty = make([]byte, 8)
   sk := []byte(strconv.Itoa(key))

   start := len(empty) - len(sk)
   var result = empty

   for i := 0; i < start; i++ {
      result[i] = 48
   }

   for i := start; i < len(empty); i++ {
      result[i] = sk[i-start]
   }
   return result
}
