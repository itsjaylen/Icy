<script lang="ts">
    import { Home, FileText, DollarSign, Phone } from 'lucide-svelte';
  
    let navLinks = [
      { name: 'Home', href: '/', icon: Home },
      { name: 'Projects', href: '/projects', icon: FileText },
      { name: 'Pricing', href: '/pricing', icon: DollarSign },
      { name: 'Contact', href: '/contact', icon: Phone }
    ];
  
    let projectsDropdownOpen = false;
    let pricingDropdownOpen = false;
    let mobileMenuOpen = false;
    let closeTimeout: ReturnType<typeof setTimeout> | undefined;
  
    const projectsSubmenu = [
      { name: 'Project A', href: '/projects/a' },
      { name: 'Project B', href: '/projects/b' },
      { name: 'Project C', href: '/projects/c' }
    ];
  
    const pricingSubmenu = [
      { name: 'Basic Plan', href: '/pricing/basic' },
      { name: 'Pro Plan', href: '/pricing/pro' },
      { name: 'Enterprise Plan', href: '/pricing/enterprise' }
    ];
  
    function toggleDropdown(type: 'projects' | 'pricing') {
      if (type === 'projects') {
        projectsDropdownOpen = !projectsDropdownOpen;
      } else {
        pricingDropdownOpen = !pricingDropdownOpen;
      }
    }
  
    function openDropdown(type: 'projects' | 'pricing') {
      clearTimeout(closeTimeout);
      if (type === 'projects') projectsDropdownOpen = true;
      if (type === 'pricing') pricingDropdownOpen = true;
    }
  
    function closeDropdown(type: 'projects' | 'pricing') {
      closeTimeout = setTimeout(() => {
        if (type === 'projects') projectsDropdownOpen = false;
        if (type === 'pricing') pricingDropdownOpen = false;
      }, 200);
    }
  </script>
  
  <nav
    class="fixed left-0 right-0 top-0 z-50 border-b border-white/10 bg-white/10 shadow-lg backdrop-blur-md dark:bg-black/20"
  >
    <div class="w-full px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-start">
        <div class="text-xl font-bold text-white">
          <a href="/">Home</a>
        </div>
        <!-- Desktop menu -->
        <div class="relative hidden items-center space-x-8 md:flex ml-auto">
          {#each navLinks as link}
            {#if link.name === 'Projects'}
              <!-- svelte-ignore a11y_no_static_element_interactions -->
              <div class="relative" on:mouseenter={() => openDropdown('projects')} on:mouseleave={() => closeDropdown('projects')}>
                <button class="font-medium text-white/80 transition duration-200 hover:text-white" on:click={() => toggleDropdown('projects')}>
                  <link.icon class="inline-block mr-2" /> {link.name}
                </button>
                {#if projectsDropdownOpen}
                  <div class="absolute top-full z-10 mt-2 w-48 rounded bg-white py-2 shadow-lg dark:bg-gray-800"
                    on:mouseenter={() => openDropdown('projects')}
                    on:mouseleave={() => closeDropdown('projects')}>
                    {#each projectsSubmenu as item}
                      <a href={item.href} class="block px-4 py-2 text-gray-700 hover:bg-gray-100 dark:text-white dark:hover:bg-gray-700">
                        {item.name}
                      </a>
                    {/each}
                  </div>
                {/if}
              </div>
            {:else if link.name === 'Pricing'}
              <!-- svelte-ignore a11y_no_static_element_interactions -->
              <div class="relative" on:mouseenter={() => openDropdown('pricing')} on:mouseleave={() => closeDropdown('pricing')}>
                <button class="font-medium text-white/80 transition duration-200 hover:text-white" on:click={() => toggleDropdown('pricing')}>
                  <link.icon class="inline-block mr-2" /> {link.name}
                </button>
                {#if pricingDropdownOpen}
                  <div class="absolute top-full z-10 mt-2 w-48 rounded bg-white py-2 shadow-lg dark:bg-gray-800"
                    on:mouseenter={() => openDropdown('pricing')}
                    on:mouseleave={() => closeDropdown('pricing')}>
                    {#each pricingSubmenu as item}
                      <a href={item.href} class="block px-4 py-2 text-gray-700 hover:bg-gray-100 dark:text-white dark:hover:bg-gray-700">
                        {item.name}
                      </a>
                    {/each}
                  </div>
                {/if}
              </div>
            {:else}
              <a href={link.href} class="font-medium text-white/80 transition duration-200 hover:text-white">
                <link.icon class="inline-block mr-2" /> {link.name}
              </a>
            {/if}
          {/each}
        </div>
  
        <!-- Mobile menu toggle -->
        <button class="text-white md:hidden" on:click={() => (mobileMenuOpen = !mobileMenuOpen)}>
          ☰
        </button>
      </div>
  
      <!-- Mobile dropdown -->
      {#if mobileMenuOpen}
        <div class="mt-2 flex flex-col space-y-2 rounded-md bg-black/70 p-4 text-white md:hidden">
          {#each navLinks as link}
            {#if link.name === 'Projects'}
              <div>
                <button class="w-full text-left font-medium" on:click={() => toggleDropdown('projects')}>
                  <link.icon class="inline-block mr-2" /> {link.name}
                </button>
                {#if projectsDropdownOpen}
                  <div class="mt-1 space-y-1 pl-4">
                    {#each projectsSubmenu as item}
                      <a href={item.href} class="block text-sm text-white/80 hover:text-white">
                        {item.name}
                      </a>
                    {/each}
                  </div>
                {/if}
              </div>
            {:else if link.name === 'Pricing'}
              <div>
                <button class="w-full text-left font-medium" on:click={() => toggleDropdown('pricing')}>
                  <link.icon class="inline-block mr-2" /> {link.name}
                </button>
                {#if pricingDropdownOpen}
                  <div class="mt-1 space-y-1 pl-4">
                    {#each pricingSubmenu as item}
                      <a href={item.href} class="block text-sm text-white/80 hover:text-white">
                        {item.name}
                      </a>
                    {/each}
                  </div>
                {/if}
              </div>
            {:else}
              <a href={link.href} class="text-white/80 hover:text-white">
                <link.icon class="inline-block mr-2" /> {link.name}
              </a>
            {/if}
          {/each}
        </div>
      {/if}
    </div>
  </nav>
  
  <style>
    nav {
      backdrop-filter: blur(12px);
    }
  </style>