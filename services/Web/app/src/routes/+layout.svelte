<script lang="ts">
	import '../app.css';
	import Footer from '$lib/components/Footer.svelte';
	import Header from '$lib/components/Header.svelte';
	import Navbar from '$lib/components/Navbar.svelte';
</script>

<div class="layout">
	<Header />
	<Navbar />
	<main>
		<slot />
	</main>
	<Footer />
</div>

<style>
	.layout {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}
	main {
		flex: 1;
		padding: 20px;
		padding-top: 80px;
		margin-top: 60px;
	}
</style>
