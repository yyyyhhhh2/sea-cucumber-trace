<script setup lang="ts"></script>

<template>
  <div class="pointer-events-none absolute inset-0 overflow-hidden">
    <div class="absolute -top-40 left-1/2 h-[520px] w-[1200px] -translate-x-1/2 rounded-full bg-hero-mesh blur-3xl" />
    <div
      class="absolute bottom-[-20%] left-[-10%] h-[420px] w-[420px] rounded-full bg-cyan-500/10 blur-3xl"
    />
    <div
      class="absolute bottom-[10%] right-[-10%] h-[380px] w-[380px] rounded-full bg-teal-500/10 blur-3xl"
    />
    <div class="noise absolute inset-0 opacity-40" />
    <svg
      class="absolute bottom-0 left-0 right-0 text-white/5"
      viewBox="0 0 1440 220"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill="currentColor"
        d="M0 120L60 105C120 90 240 60 360 52.5C480 45 600 60 720 75C840 90 960 105 1080 97.5C1200 90 1320 60 1380 45L1440 30V220H0V120Z"
      />
    </svg>
  </div>
</template>
