<script setup lang="ts"></script>

<template>
  <div class="pointer-events-none absolute inset-0 overflow-hidden">
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(10,132,180,0.18),transparent_38%),linear-gradient(180deg,rgba(6,34,52,0.02),rgba(2,6,23,0.86))]" />
    <div class="coast-grid absolute inset-0 opacity-[0.08]" />

    <div class="absolute left-1/2 top-[-18rem] h-[42rem] w-[78rem] -translate-x-1/2 rounded-full bg-hero-mesh opacity-80 blur-3xl animate-drift" />
    <div class="absolute left-[8%] top-[18%] h-72 w-72 rounded-full bg-cyan-400/10 blur-3xl animate-sway" />
    <div class="absolute right-[6%] top-[26%] h-80 w-80 rounded-full bg-teal-400/10 blur-3xl animate-drift" />
    <div class="absolute bottom-[-8rem] left-[12%] h-64 w-64 rounded-full bg-sky-300/10 blur-3xl animate-float" />

    <div class="absolute inset-x-0 top-[10%] flex justify-center opacity-50">
      <div class="h-[28rem] w-[28rem] rounded-full border border-white/[0.04] animate-spin-slow" />
      <div class="absolute mt-14 h-[17rem] w-[17rem] rounded-full border border-cyan-300/[0.07] animate-spin-slow [animation-direction:reverse]" />
    </div>

    <svg
      class="absolute inset-x-0 bottom-[17vh] h-36 w-full text-slate-200/[0.05]"
      viewBox="0 0 1440 220"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      preserveAspectRatio="none"
    >
      <path
        fill="currentColor"
        d="M0 183H84V109H118V182H152V142H190V183H230V92H280V183H322V123H364V183H426V153H462V183H530V74H589V183H626V128H670V183H720V99H764V183H805V149H840V183H899V118H944V183H1004V94H1055V183H1110V131H1144V183H1206V108H1260V183H1300V146H1335V183H1382V122H1440V220H0V183Z"
      />
    </svg>

    <div class="absolute inset-x-0 bottom-[10vh] h-[18vh] overflow-hidden">
      <svg
        class="absolute inset-x-0 bottom-10 h-40 w-full text-white/[0.05]"
        viewBox="0 0 1440 320"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        preserveAspectRatio="none"
      >
        <path
          fill="currentColor"
          d="M0 224L40 213.3C80 203 160 181 240 186.7C320 192 400 224 480 234.7C560 245 640 235 720 213.3C800 192 880 160 960 154.7C1040 149 1120 171 1200 176C1280 181 1360 171 1400 165.3L1440 160V320H1400C1360 320 1280 320 1200 320C1120 320 1040 320 960 320C880 320 800 320 720 320C640 320 560 320 480 320C400 320 320 320 240 320C160 320 80 320 40 320H0V224Z"
        />
      </svg>
      <svg
        class="absolute inset-x-0 bottom-0 h-44 w-full text-cyan-200/[0.06]"
        viewBox="0 0 1440 320"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        preserveAspectRatio="none"
      >
        <path
          fill="currentColor"
          d="M0 288L48 266.7C96 245 192 203 288 202.7C384 203 480 245 576 250.7C672 256 768 224 864 192C960 160 1056 128 1152 133.3C1248 139 1344 181 1392 202.7L1440 224V320H1392C1344 320 1248 320 1152 320C1056 320 960 320 864 320C768 320 672 320 576 320C480 320 384 320 288 320C192 320 96 320 48 320H0V288Z"
        />
      </svg>
    </div>

    <div class="noise absolute inset-0 opacity-35" />
  </div>
</template>
