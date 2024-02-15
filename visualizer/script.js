import * as THREE from 'three'

const scene = new THREE.Scene()
const camera = new THREE.PerspectiveCamera( 120, window.innerWidth / window.innerHeight, 1, 1000 )

const renderer = new THREE.WebGLRenderer()
renderer.setSize( window.innerWidth, window.innerHeight )
document.body.appendChild( renderer.domElement )

const geometry = new THREE.CylinderGeometry( 20, 20, 40, 32 ); 
// const T01 = new THREE.Mesh( geometry, material ); scene.add( T01 );
// const T12 = new THREE.Mesh( geometry, material ); scene.add( T12 );
// const T23 = new THREE.Mesh( geometry, material ); scene.add( T23 );
// const T34 = new THREE.Mesh( geometry, material ); scene.add( T34 );
// const T45 = new THREE.Mesh( geometry, material ); scene.add( T45 );
// const T56 = new THREE.Mesh( geometry, material ); scene.add( T56 );

let Tarr = []

for(let i = 0; i < 6; i++) {
  const material = new THREE.MeshBasicMaterial( {color: Math.floor(Math.random()*16777215), wireframe: true} )
  Tarr[i] = new THREE.Mesh( geometry, material )
  scene.add(Tarr[i])
}

let rawMtxArr = []
async function updateTransformations() {
for (let i = 0; i < 6; i++) {
  await fetch("http://localhost:8080/matrix/" + i + (i+1)).then(d => d.text()).then(raw => {
    // Filter the response
    const rawMtx = raw.split(" ")
    var filtered = []
    for (let i = 0; i < rawMtx.length; i++) {
      let d = rawMtx[i]
      if (typeof Number(d) != "number" || d === "" || d === "\n") {
        continue
      }
      if (d.startsWith("-0")) {
        d = "0"
      }
      filtered.push(d)
    }
    filtered.pop()

    // Initialize matrix
    rawMtxArr[i] = new THREE.Matrix4()
    rawMtxArr[i].fromArray(filtered).transpose()
    console.table(rawMtxArr[0].elements)
    console.table(Tarr[0].matrix.elements)
  })
}}

camera.position.z = 550
camera.position.y = 400
camera.position.x = 200
camera.rotation.y = -.5

function animate() {
  requestAnimationFrame(animate)
if (rawMtxArr.length == 6) 
  // for(let i = 0; i < 6; i++) {
  //   Tarr[i].applyMatrix4(rawMtxArr[i])
  //   Tarr[i].updateMatrix()
  // }
  renderer.render(scene, camera)
}

updateTransformations()
animate()
