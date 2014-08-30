weather
=======

A webserver that runs on a Raspberry Pi and lets you:
* control Philips Hue lights
* control the volume of the headphone jack
* connect to a streaming audio source securely

# Setup

Depends on my [volume package](https://github.com/bklimt/volume) and my [Hue package](https://github.com/bklimt/hue), so make sure to install the dependencies for those. Specifically, you'll need to make sure you install the ALSA asound development libraries.

You'll need to make sure your Hue is configured correctly. Build the `hue` command-line tool.

    go install github.com/bklimt/hue/...

Then press the link button on your Hue router and use the command-line tool to register the app.

    $GOPATH/bin/hue --hue_ip="YOUR.ROUTER.IP.ADDRESS" --register

You'll also need to set up mysql.

    // instructions coming soon

Once that it complete, you can start the `weather` server:

    $GOPATH/bin/humid --hue_ip="YOUR.ROUTER.IP.ADDRESS"

While `weather` is running, you can access it using the IP address of the Raspberry Pi.

