<script>
	export let config = null
	config = {
		...(config || {})
	}
</script>

<template lang="pug">
.form-item
  label(for='bucket') Bucket
  input.input.mt-1(type='text' name='bucket' bind:value='{config.bucket}')
.form-item
  label(for='location') Location
  input.input.mt-1(type='text' name='workdir' bind:value='{config.location}')
.form-item
  label(for='endpoint') Endpoint
  input.input.mt-1(type='text' name='endpoint' bind:value='{config.endpoint}')
.form-item
  label(for='accessKey') AWS Access Key ID
  input.input.mt-1(type='text' name='accessKey' bind:value='{config.accessKeyId}')
.form-item
  label(for='secretAccessKey') AWS Secret Access Key
  input.input.mt-1(type='password' name='secretAccessKey' bind:value='{config.secretAccessKeyId}')
</template>
