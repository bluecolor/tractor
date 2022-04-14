<script>
	import _ from 'lodash';
	export let items = [];
	export let valueField = 'id';
	export let textField = 'name';
	export let loading = false;
	export let readonly = true;
	export let value = null;

	let open = false;

	$: _items = items.map((item) => {
		if (_.isObject(item)) {
			return {
				value: item[valueField],
				text: item[textField]
			};
		}
		return {
			value: item,
			text: item
		};
	});

	function toggle() {
		open = !open;
	}
</script>

<template lang="pug">
.rich-select.w-full.open(tabindex='0' on:click='{toggle}')
  +if('readonly')
    span.current.pl-4.text-gray-600 Select option
    +else
      input(type='text', placeholder='Search')
  +if('loading')
    <img src="src/assets/icons/loading.svg" class="loading" alt="loading"/>
  +if('open')
    .rich-select-dropdown.w-full
      ul.list
        li.option(data-value='Nothing') Nothing
        li.option(data-value='1') Some option
        li.option.selected.focus(data-value='2') Another option
        li.option.disabled(data-value='3') A disabled option
        li.option(data-value='4') Potato

</template>

<style lang="scss">
	.rich-select.open {
		@apply ring-blue-200
        ring-offset-blue-200
        ring-2
        ring-offset-1
        outline-none
        border-0;
	}
	.rich-select {
		.loading {
			position: absolute;
			height: 1.5rem;
			width: 1.5rem;
			top: 0.4rem;
			right: 0.5rem;
		}
		-webkit-tap-highlight-color: rgba(0, 0, 0, 0);
		background-color: #fff;
		border-radius: 5px;
		border: solid 1px #e8e8e8;
		box-sizing: border-box;
		clear: both;
		cursor: pointer;
		display: block;
		float: left;
		font-family: inherit;
		font-size: 14px;
		font-weight: normal;
		height: 38px;
		line-height: 36px;
		outline: none;
		position: relative;
		text-align: left !important;
		transition: all 0.2s ease-in-out;
		user-select: none;
		white-space: nowrap;
		input {
			@apply block rounded-md border-0 w-full outline-none;
			height: calc(100% - 2px);
			width: calc(100% - 2px);
			&:focus {
				--tw-ring-shadow: 0;
			}
		}
		.rich-select-dropdown {
			margin-top: 8px;
			background-color: #fff;
			border-radius: 5px;
			box-shadow: 0 0 0 1px rgba(68, 68, 68, 0.11);
			pointer-events: none;
			position: absolute;
			top: 100%;
			left: 0;
			transform-origin: 50% 0;
			transform: scale(0.75) translateY(19px);
			transition: all 0.2s cubic-bezier(0.5, 0, 0, 1.25), opacity 0.15s ease-out;
			z-index: 9;
			opacity: 1;
			pointer-events: auto;
			transform: scale(1) translateY(0);
		}
		.list {
			border-radius: 5px;
			box-sizing: border-box;
			overflow: hidden;
			padding: 0;
			max-height: 210px;
			overflow-y: auto;
			li.option {
				padding-left: 18px;
				padding-right: 18px;
				&:hover {
					background-color: #f5f5f5;
				}
			}
		}
	}
</style>
