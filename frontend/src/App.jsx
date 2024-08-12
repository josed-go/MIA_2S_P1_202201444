import Editor from '@monaco-editor/react'
import FileInput from './components/FileInput'
function App() {

  return (
    <div className="h-screen flex flex-col text-center justify-center">
      <section className='flex flex-col text-center justify-center'>
        <h1 className='m-4 text-2xl'>Entrada</h1>
        <div className='flex justify-center'>
          <Editor className='rounded-md'
              height="15vh" 
              width="50%"
              theme='vs-dark'
              defaultLanguage='cpp'
              defaultValue=''
              options={{
                scrollBeyondLastLine:false,
                fontSize:"20px"
            }}
          />
        </div>
        <FileInput />
        <button className='my-6 mx-auto p-2 rounded-md bg-btn w-1/12 text-xl font-bold text-white'>Ejecutar</button>
      </section>
      <section className='flex flex-col text-center justify-center'>
        <h1 className='m-4 text-2xl'>Salida</h1>
        <div className='flex justify-center'>
          <Editor className='rounded-md'
              height="15vh" 
              width="50%"
              theme='vs-dark'
              defaultLanguage='cpp'
              defaultValue=''
              options={{
                scrollBeyondLastLine:false,
                fontSize:"20px",
                readOnly: true
              }}
          />
        </div>
      </section>
    </div>
  )
}

export default App
