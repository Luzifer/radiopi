# Luzifer / radiopi

Imaging this: You own a car, not the newest generation but already a few years old. Your car has a AUX input to the sound system and got a quite decent sound system... The local radio stations are full of music you don't like and burning CDs is just something you don't want to do every time.

I am in exactly this situation. Until now I'm connecting my smartphone to the car sound system and playing streams using different apps. All those apps are not really reliable, the streams disconnect often and also my battery gets leeched. But also I'm working in an industrial section where I'm automating things. Including setups of Raspberry-PIs. So this idea got born.

This project is intended to run on a 2nd generation Raspberry-PI and is set up using [Ansible](http://www.ansible.com/home). Also a really tiny [Go](https://golang.org/) daemon is deployed on that device. The daemon is running on start of the device and the device is set up to keep connected to a Huawei WiFi stick. Now put everyhing together and you got an internet radio for your car powered using just an USB car adapter automatically starting up and connecting to the last station you played.

Also it will have a small web interface you can use to tune in to another web radio station...
