non hit specular = 0x3498db
non hit color = 0x8e44ad
non hit emissive = 0x34495e
hit specular = 0xf39c12
hit color = 0xd35400
hit emissive = 0xc0392b

init camera () =
  far          = 10000
  near         = 0.1
  view_angle   = 45
  aspect       = window.inner width / window.inner height
  c            = @new THREE.Perspective camera(view_angle, aspect, near, far)
  c.position.z = 300
  c

add lights (scene) =
  light = @new THREE.AmbientLight( 0x222222 )
  scene.add(light)
  directional light = @new THREE.Directional light 0xffffff, 0.2
  directional light.position.set 0 0 1000
  directional light.target.position.copy( scene.position )
  scene.add(directional light)

init scene () =
  scene = @new THREE.Scene
  scene

init controls (camera) =
  c = @new THREE.Trackball controls(camera)

  c.rotateSpeed = 1.0
  c.zoomSpeed = 1.2
  c.panSpeed = 0.8

  c.noZoom = false
  c.noPan = false

  c.staticMoving = true
  c.dynamicDampingFactor = 0.3

  c.keys = [ 65, 83, 68 ]

  c.addEventListener( 'change', render )
  c

init renderer () =
  renderer = @new THREE.WebGL renderer(antialias: true)
  renderer.set size (window.innerWidth, window.innerHeight)
  renderer

cantors pairing (a,b) =
  0.5 * (a + b) * (a + b + 1) + b

init grid (scene) =
  grid = {}
  for (x = 0, x < window.inner width / 3, x := x + 10)
    for (y = 0, y < window.inner height / 3, y := y + 10)
      cube material = @new THREE.Mesh phong material(
        specular: (non hit specular)
        color: (non hit color)
        emissive: (non hit emissive)
        transparent: true
        opacity: 0.7
        shininess: 200
        shading: THREE.SmoothShading
      )
      cube = @new THREE.Mesh(
        @new THREE.Cube geometry(6, 6, 2)
        cube material
      )
      cube x = x - window.inner width / 6
      cube y = window.inner height / 6 - y
      cube.position.set(cube x, cube y, 0)
      cube.matrixAutoUpdate = false
      cube.updateMatrix()
      scene.add (cube)

      grid.(cantors pairing (x/10, y/10)) = {
        cube = cube
      }

  grid


animate () =
  request animation frame (animate)
  render next (generation)
  controls.update()
  render ()

render () =
  renderer.render(scene, camera)

on window resize () =
  camera.aspect = window.inner width / window.inner height
  camera.updateProjectionMatrix()
  renderer.setSize( window.innerWidth, window.innerHeight )
  controls.handleResize()
  render()

render next (generation) =
  for @(key) in (grid)
    hit = false
    cube = grid.(key).cube

    for each @(point) in (generation)
      if (Number(key) == cantors pairing (point.Col, point.Row))
        cube.material.specular.set hex (hit specular)
        cube.material.color.set hex (hit specular)
        cube.material.emissive.set hex (hit emissive)
        hit := true

    if (! hit)
      cube.material.specular.set hex (non hit specular)
      cube.material.color.set hex (non hit color)
      cube.material.emissive.set hex (non hit emissive)

    cube.geometry.colorsNeedUpdate = true

camera   = init camera ()
scene    = init scene ()
controls = init controls (camera)
renderer = init renderer ()
add lights (scene)
grid = init grid (scene)
generation = []

container = document.get element by id "container"
container.append child (renderer.dom element)

animate ()

ws = @new Web socket "ws://#(window.location.host)/go-ws"

window.addEventListener( 'resize', onWindowResize, false )

ws.onmessage (event) =
  generation := JSON.parse(event.data)
