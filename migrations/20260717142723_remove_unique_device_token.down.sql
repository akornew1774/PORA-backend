ALTER TABLE devices
ADD CONSTRAINT devices_device_token_key UNIQUE (device_token);