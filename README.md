# An example of a storage from a hashicorp vault
Bulbous storage from hahshicorp vault. All layers are thrown into each other's constructors and implement the same Storage interface

First Layer: Physical storage(this object works with the database directly)
Second Layer: Cache storage(lru cacke)
Third layer: Encrypted storage(aes-gcm alg)
Fourth layer: Barrier with prefix(you can make an object that constantly adds a prefix to the keys. It is necessary to delimit areas in the storage)
