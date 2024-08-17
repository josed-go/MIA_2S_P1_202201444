import Editor from '@monaco-editor/react'
import FileInput from './components/FileInput'
import { useRef, useState} from 'react'

function App() {
  const editorRef = useRef(null)
  const consolaRef = useRef(null)

  const [ entradaFile, setEntradaFile ] = useState("")

  const handleEditor = (editor, id) => {
    if(id == "editor" ) {
        editorRef.current = editor
    }else if(id == "consola") {
        consolaRef.current = editor
    }
  }

  const extractPath = (command) => {
    const pathFlag = "-path=";
    const startIndex = command.indexOf(pathFlag);
  
    if (startIndex !== -1) {
      return command.substring(startIndex + pathFlag.length).split(' ')[0];
    }
  
    return null
  }

  const confirmarRmdisk = (entrada) => {
    const lineas = entrada.split('\n');

    // Almacena las líneas filtradas
    let lineasFiltradas = [];

    // Itera sobre cada línea
    for (let i = 0; i < lineas.length; i++) {
      const linea = lineas[i]

      if (linea.includes('rmdisk')) {
        const ruta = extractPath(linea)
        const confirmar = window.confirm(`Comando en linea ${i+1} eliminara la ruta:\n${ruta}.\nQuieres continuar?`)

        if(confirmar) {
          lineasFiltradas.push(linea)
        } 
        continue
      }
      lineasFiltradas.push(linea);
    }

    // Une las líneas filtradas y devuelve el nuevo string
    return lineasFiltradas.join('\n');
  }

  const analizar = async () => {
    var entrada = editorRef.current.getValue()

    entrada = entrada +"\n"+ entradaFile
    
    const entradaFiltrada = confirmarRmdisk(entrada)

    const response = await fetch('http://localhost:3000/analizar', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ entrada: entradaFiltrada }),
    })

    const data = await response.json()
    consolaRef.current.setValue(data.response)  
  }

  return (
    <div className="h-screen flex flex-col text-center justify-center">
      <section className='flex flex-col text-center justify-center'>
        <h1 className='m-4 text-2xl'>Entrada</h1>
        <div className='flex justify-center'>
          <Editor className='rounded-md'
              height="25vh" 
              width="55%"
              theme='vs-dark'
              defaultLanguage='cpp'
              defaultValue=''
              options={{
                scrollBeyondLastLine:false,
                fontSize:"16px"
              }}
              onMount={(editor) => handleEditor(editor, "editor")}
          />
        </div>
        <FileInput texto={setEntradaFile} />
        <button className='my-6 mx-auto p-2 rounded-md bg-btn w-1/12 text-xl font-bold text-white hover:bg-btn-osc'
          onClick={analizar}
        >
          Ejecutar
        </button>
      </section>
      <section className='flex flex-col text-center justify-center'>
        <h1 className='m-4 text-2xl'>Salida</h1>
        <div className='flex justify-center'>
          <Editor className='rounded-md'
              height="25vh" 
              width="55%"
              theme='vs-dark'
              defaultLanguage='cpp'
              defaultValue=''
              options={{
                scrollBeyondLastLine:false,
                fontSize:"16px",
                readOnly: true
              }}
              onMount={(editor) => handleEditor(editor, "consola")}
          />
        </div>
      </section>
    </div>
  )
}

export default App
