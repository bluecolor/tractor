<script>
	export let state = {};
</script>

<template lang="pug">
  .form-item
    label(for='host-name') Host
    input.mt-1(type='text' id="host-name" name='host' autocomplete='mysql-host-name' bind:value='{state.host}')

  .form-item
    label(for='port') Port
    input.mt-1(type='number' id="port" name='port' autocomplete='mysql-port' bind:value='{state.port}')

  .form-item
    label(for='database') Database
    input.mt-1(type='text' id="database" name='database' autocomplete='mysql-database' bind:value='{state.database}')

  .form-item
    label(for='username') Username
    input.mt-1(type='text' id="username" name='username' autocomplete='mysql-username' bind:value='{state.username}')

  .form-item
    label(for='password') Password
    input.mt-1(type='password' id="password" name='password' autocomplete='mysql-password' bind:value='{state.password}')

</template>
