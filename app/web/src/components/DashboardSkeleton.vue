<template>
    <!-- Stats skeleton -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <div v-for="index in 4" :key="index" class="rounded-xl border border-base-300/60 bg-base-100 p-5">
            <div class="flex items-start justify-between">
                <div class="space-y-3">
                    <div class="skeleton h-4 w-24"></div>
                    <div class="skeleton h-8 w-32"></div>
                    <div class="skeleton h-3 w-24"></div>
                </div>

                <div class="skeleton h-10 w-10 rounded-lg"></div>
            </div>
        </div>
    </div>

    <div class="grid grid-cols-1 gap-5">
        <div class="rounded-xl border border-base-300/60 bg-base-100 p-5 ">
            <div class="mb-5 space-y-2">
                <div class="skeleton h-5 w-32"></div>
                <div class="skeleton h-3 w-56"></div>
            </div>

            <div class="skeleton h-80 w-full"></div>
        </div>

    </div>

    <div class="grid grid-cols-1 gap-5 xl:grid-cols-3">

        <div class="rounded-xl border border-base-300/60 bg-base-100 p-5">
            <div class="mb-5 space-y-2">
                <div class="skeleton h-5 w-40"></div>
                <div class="skeleton h-3 w-48"></div>
            </div>

            <div class="skeleton mx-auto h-52 w-52 rounded-full"></div>
        </div>


        <div class="rounded-xl border border-base-300/60 bg-base-100 p-5 xl:col-span-2">
            <div class="mb-5 flex items-center justify-between">
                <div class="space-y-2">
                    <div class="skeleton h-5 w-40"></div>
                    <div class="skeleton h-3 w-52"></div>
                </div>

                <div class="skeleton h-8 w-20"></div>
            </div>

            <div class="space-y-1">
                <div v-for="index in 5" :key="index"
                    class="flex items-center gap-4 border-b border-base-300/40 py-4 last:border-0">
                    <div class="skeleton h-10 w-10 shrink-0 rounded-lg"></div>

                    <div class="min-w-0 flex-1 space-y-2">
                        <div class="skeleton h-4 w-28"></div>
                        <div class="skeleton h-3 w-20"></div>
                    </div>

                    <div class="skeleton h-4 w-20"></div>
                </div>
            </div>
        </div>
    </div>

</template>
