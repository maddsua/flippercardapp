<template>
	<div class="form">
		<div class="wrapper">
			<slot>
				[Form]
			</slot>
		</div>
	</div>
</template>

<style lang="scss" scoped>
	.form {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		flex-grow: 1;
		
		.wrapper {
			display: flex;
			flex-direction: column;
			gap: 2rem;
			width: 100%;
		}
	}
</style>
