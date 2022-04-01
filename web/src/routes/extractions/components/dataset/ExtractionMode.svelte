<script>
	export let type = 'source';
	export let value;
	let modes = [];
	console.log(type);
	switch (type) {
		case 'target':
			modes = [
				{
					label: 'Append',
					value: 'append'
				},
				{
					label: 'Replace',
					value: 'replace'
				},
				{
					label: 'Merge',
					value: 'merge'
				}
			];
			break;
		case 'source':
			modes = [
				{
					label: 'Full',
					value: 'full'
				},
				{
					label: 'Incremental',
					value: 'incremental'
				}
			];
			break;
	}
	function onModeSelect(m) {
		value = m.value;
	}
</script>

<template lang="pug">
.form-item
  label(for="table") Mode
  .flex
    +each('modes as m')
      div(on:click='{onModeSelect(m)}')
        label.inline-flex.items-center.mr-2()
          input.form-radio.h-4.w-4.text-blue-400(type='radio' checked='{value==m.value}'  class='focus:ring-2 focus:ring-offset-1 focus:border-blue-300 focus:ring focus:ring-offset-0 focus:ring-blue-200')
          span.ml-2.text-gray-400 {m.label}

</template>
