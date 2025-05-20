<template>
  <div class="app">
    <NavBar />
    <Hero />
    <About />
    <Skills />
    <Experience />
    <Projects />
    <Certificates />
    <Contact />
    <Footer />
  </div>
</template>

<script setup>
import NavBar from './components/NavBar.vue'
import Hero from './components/Hero.vue'
import About from './components/About.vue'
import Skills from './components/Skills.vue'
import Experience from './components/Experience.vue'
import Projects from './components/Projects.vue'
import Certificates from './components/Certificates.vue'
import Contact from './components/Contact.vue'
import Footer from './components/Footer.vue'
</script>

<style>
:root {
  --primary-color: #1E3A8A;
  --secondary-color: #3B82F6;
  --accent-color: #06B6D4;
  --text-color: #374151;
  --light-color: #F3F4F6;
  --dark-color: #1F2937;
  --transition-speed: 0.3s;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: "PingFang SC", "Microsoft YaHei", sans-serif;
  color: var(--text-color);
  background-color: #f8f9fa;
  overflow-x: hidden;
}

.app {
  width: 100%;
  min-height: 100vh;
}

section {
  padding: 80px 0;
  overflow: hidden;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

h1, h2, h3, h4, h5, h6 {
  font-family: "Source Han Sans CN", "Microsoft YaHei", sans-serif;
  color: var(--primary-color);
  margin-bottom: 1rem;
}

a {
  color: var(--secondary-color);
  text-decoration: none;
  transition: all var(--transition-speed);
}

a:hover {
  color: var(--accent-color);
}

.section-title {
  text-align: center;
  margin-bottom: 3rem;
  position: relative;
}

.section-title h2 {
  font-size: 2.5rem;
  display: inline-block;
  padding-bottom: 10px;
  position: relative;
}

.section-title h2:after {
  content: '';
  position: absolute;
  display: block;
  width: 50px;
  height: 3px;
  background: var(--accent-color);
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}

@media (max-width: 768px) {
  section {
    padding: 60px 0;
  }
  
  .section-title h2 {
    font-size: 2rem;
  }
}
</style> 