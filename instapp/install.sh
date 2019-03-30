#!/bin/sh -e
if [ $# -eq 0 ] ; then
    echo 'argument expected like install.sh app_name'
    exit 1
fi

APP=$1
service $APP stop
rm -r /usr/$APP 2>/dev/null
rm /var/log/$APP 2>/dev/null
rm /usr/bin/configuration.yaml 2>/dev/null
rm /etc/init.d/$APP.sh 2>/dev/null
mkdir /usr/$APP
mkdir /usr/$APP/png
echo "htmldir: /usr/${APP}/html/" >> /usr/bin/configuration.yaml
echo "htmlfile: /usr/${APP}/html/index.html" >> /usr/bin/configuration.yaml
echo "appdir:  /usr/${APP}/" >> /usr/bin/configuration.yaml
echo "dbname:  /usr/${APP}/${APP}.db" >> /usr/bin/configuration.yaml
echo "logfile: /var/log/${APP}.log" >> /usr/bin/configuration.yaml
echo "pngdir: /usr/${APP}/png/" >> /usr/bin/configuration.yaml
echo "addr:    :443" >> /usr/bin/configuration.yaml
echo "key:     /etc/ssl/goyav/goyav.com.key" >> /usr/bin/configuration.yaml
echo "crt:     /etc/ssl/goyav/goyav.com.crt" >> /usr/bin/configuration.yaml
echo "pem:     /etc/ssl/goyav/GandiStandardSSLCA2.pem" >> /usr/bin/configuration.yaml
echo "host:     mail.gandi.net:25" >> /usr/bin/configuration.yaml
echo "user:     admin@goyav.com" >> /usr/bin/configuration.yaml
echo "password: HpNKUsNN@27031968" >> /usr/bin/configuration.yaml

cp -r html /usr/$APP
cp -r goyav /etc/ssl
cp $APP /usr/bin
update-rc.d -f $APP.sh remove
perl -pi.back -e "s/MYSERVICE/${APP}/g;" service.sh
cp service.sh /etc/init.d/$APP.sh
chmod 0755 /etc/init.d/$APP.sh
chmod 0755 /etc/ssl/goyav
chmod 0755 /usr/bin/configuration.yaml
systemctl daemon-reload
update-rc.d $APP.sh defaults
systemctl daemon-reload
service $APP start
service $APP stop
chmod 0755 /usr/$APP/$APP.db
service $APP start
