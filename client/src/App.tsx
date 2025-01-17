import { Container, Stack } from '@chakra-ui/react'
import Navbar from './Components/Navbar'
import TodoForm from './Components/TodoForm'
import TodoList from './Components/TodoList'

//export const BASE_URL = import.meta.env.MODE == "development"   ?  "http://localhost:5000/api/todos/" : 

export const BASE_URL = "/api/todos/"

function App() {
  return (
    <Stack h="100vh">
     <Navbar />
     <TodoForm />
     <TodoList />
     <Container />
    </Stack>
  )
}

export default App
